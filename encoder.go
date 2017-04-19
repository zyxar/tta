package tta

import (
	"fmt"
	"io"
	"os"
)

type Encoder struct {
	codecs    [maxNCH]codec // 1 per channel
	channels  int           // number of channels/codecs
	data      [8]byte       // codec initialization data
	fifo      fifo
	seekTable []uint64 // the playing position table
	format    uint32   // tta data format
	rate      uint32   // bitrate (kbps)
	offset    uint64   // data start position (header size, bytes)
	frames    uint32   // total count of frames
	depth     uint32   // bytes per sample
	flenStd   uint32   // default frame length in samples
	flenLast  uint32   // last frame length in samples
	flen      uint32   // current frame length in samples
	fnum      uint32   // currently playing frame index
	fpos      uint32   // the current position in frame
	shiftBits uint32   // packing int to pcm
}

func Compress(infile, outfile io.ReadWriteSeeker, passwd string, cb Callback) (err error) {
	waveHdr := WaveHeader{}
	var dataSize uint32
	if dataSize, err = waveHdr.Read(infile); err != nil {
		err = errRead
		return
	} else if dataSize >= 0x7FFFFFFF {
		err = fmt.Errorf("incorrect data size info in wav file: %x", dataSize)
		return
	}
	if (waveHdr.chunkId != riffSign) ||
		(waveHdr.format != waveSign) ||
		(waveHdr.numChannels == 0) ||
		(waveHdr.numChannels > maxNCH) ||
		(waveHdr.bitsPerSample == 0) ||
		(waveHdr.bitsPerSample > maxBPS) {
		err = errFormat
		return
	}
	encoder := NewEncoder(outfile)
	smpSize := uint32(waveHdr.numChannels * ((waveHdr.bitsPerSample + 7) / 8))
	info := Info{
		nch:     uint32(waveHdr.numChannels),
		bps:     uint32(waveHdr.bitsPerSample),
		sps:     waveHdr.sampleRate,
		format:  formatSimple,
		samples: dataSize / smpSize,
	}
	if len(passwd) > 0 {
		encoder.SetPassword(passwd)
		info.format = formatEncrypted
	}
	bufSize := pcmBufferLength * smpSize
	buffer := make([]byte, bufSize)
	if err = encoder.SetInfo(&info, 0); err != nil {
		return
	}
	var readLen int
	for dataSize > 0 {
		if bufSize >= dataSize {
			bufSize = dataSize
		}
		if readLen, err = infile.Read(buffer[:bufSize]); err != nil || readLen != int(bufSize) {
			err = errRead
			return
		}
		encoder.ProcessStream(buffer[:bufSize], cb)
		dataSize -= bufSize
	}
	encoder.Close()
	return
}

func NewEncoder(iocb io.ReadWriteSeeker) *Encoder {
	enc := Encoder{}
	enc.fifo.io = iocb
	return &enc
}

func (e *Encoder) ProcessStream(in []byte, cb Callback) {
	if len(in) == 0 {
		return
	}
	var res, curr, next, tmp int32
	next = readBuffer(in, e.depth)
	in = in[e.depth:]
	tmp = next << e.shiftBits
	i := 0
	index := 0
	for {
		curr = next
		if index < len(in) {
			next = readBuffer(in[index:], e.depth)
			tmp = next << e.shiftBits
		} else {
			next = 0
			tmp = 0
		}
		index += int(e.depth)
		// transform data
		if e.channels > 1 {
			if i < e.channels-1 {
				res = next - curr
				curr = res
			} else {
				curr -= res / 2
			}
		}
		// compress stage 1: fixed order 1 prediction
		tmp = curr
		curr -= ((e.codecs[i].prev * ((1 << 5) - 1)) >> 5)
		e.codecs[i].prev = tmp
		// compress stage 2: adaptive hybrid filter
		e.codecs[i].filter.Encode(&curr)
		e.fifo.putValue(&e.codecs[i].adapter, curr)
		if i < e.channels-1 {
			i++
		} else {
			i = 0
			e.fpos++
		}
		if e.fpos == e.flen {
			e.fifo.flushBitCache()
			e.seekTable[e.fnum] = uint64(e.fifo.count)
			e.fnum++
			// update dynamic info
			e.rate = (e.fifo.count << 3) / 1070
			if cb != nil {
				cb(e.rate, e.fnum, e.frames)
			}
			e.frameInit(e.fnum)
		}
		if index >= int(e.depth)+len(in) {
			break
		}
	}
}

func (e *Encoder) ProcessFrame(in []byte) {
	if len(in) == 0 {
		return
	}
	var res, curr, next, tmp int32
	next = readBuffer(in, e.depth)
	in = in[e.depth:]
	tmp = next << e.shiftBits
	i := 0
	index := 0
	for {
		curr = next
		if index < len(in) {
			next = readBuffer(in[index:], e.depth)
			tmp = next << e.shiftBits
		} else {
			next = 0
			tmp = 0
		}
		index += int(e.depth)
		// transform data
		if e.channels > 1 {
			if i < e.channels-1 {
				res = next - curr
				curr = res
			} else {
				curr -= res / 2
			}
		}
		// compress stage 1: fixed order 1 prediction
		tmp = curr
		curr -= ((e.codecs[i].prev * ((1 << 5) - 1)) >> 5)
		e.codecs[i].prev = tmp
		// compress stage 2: adaptive hybrid filter
		e.codecs[i].filter.Encode(&curr)
		e.fifo.putValue(&e.codecs[i].adapter, curr)
		if i < e.channels-1 {
			i++
		} else {
			i = 0
			e.fpos++
		}
		if e.fpos == e.flen {
			e.fifo.flushBitCache()
			// update dynamic info
			e.rate = (e.fifo.count << 3) / 1070
			break
		}
		if index >= int(e.depth)+len(in) {
			break
		}
	}
}

func (e *Encoder) writeSeekTable() (err error) {
	if e.seekTable == nil {
		return
	}
	if _, err = e.fifo.io.Seek(int64(e.offset), os.SEEK_SET); err != nil {
		return
	}
	e.fifo.writeStart()
	e.fifo.reset()
	for i := uint32(0); i < e.frames; i++ {
		e.fifo.writeUint32(uint32(e.seekTable[i] & 0xFFFFFFFF))
	}
	e.fifo.writeCrc32()
	e.fifo.writeDone()
	return
}

func (e *Encoder) SetPassword(pass string) {
	e.data = computeKeyDigits(convertPassword(pass))
}

func (e *Encoder) frameInit(frame uint32) (err error) {
	if frame >= e.frames {
		return
	}
	shift := shifts[e.depth-1]
	e.fnum = frame
	if e.fnum == e.frames-1 {
		e.flen = e.flenLast
	} else {
		e.flen = e.flenStd
	}
	// init entropy encoder
	for i := 0; i < e.channels; i++ {
		e.codecs[i].filter = NewCompatibleFilter(e.data, shift)
		e.codecs[i].adapter.init(10, 10)
		e.codecs[i].prev = 0
	}
	e.fpos = 0
	e.fifo.reset()
	return
}

func (e *Encoder) frameReset(frame uint32, iocb io.ReadWriteSeeker) {
	e.fifo.io = iocb
	e.fifo.readStart()
	e.frameInit(frame)
}

func (e *Encoder) WriteHeader(info *Info) (size uint32, err error) {
	e.fifo.reset()
	// write TTA1 signature
	if err = e.fifo.writeByte('T'); err != nil {
		return
	}
	if err = e.fifo.writeByte('T'); err != nil {
		return
	}
	if err = e.fifo.writeByte('A'); err != nil {
		return
	}
	if err = e.fifo.writeByte('1'); err != nil {
		return
	}
	if err = e.fifo.writeUint16(uint16(info.format)); err != nil {
		return
	}
	if err = e.fifo.writeUint16(uint16(info.nch)); err != nil {
		return
	}
	if err = e.fifo.writeUint16(uint16(info.bps)); err != nil {
		return
	}
	if err = e.fifo.writeUint32(info.sps); err != nil {
		return
	}
	if err = e.fifo.writeUint32(info.samples); err != nil {
		return
	}
	if err = e.fifo.writeCrc32(); err != nil {
		return
	}
	size = 22
	return

}

func (e *Encoder) SetInfo(info *Info, pos int64) (err error) {
	if info.format > 2 ||
		info.bps < minBPS ||
		info.bps > maxBPS ||
		info.nch > maxNCH {
		return errFormat
	}
	// set start position if required
	if pos != 0 {
		if _, err = e.fifo.io.Seek(int64(pos), os.SEEK_SET); err != nil {
			err = errSeek
			return
		}
	}
	e.fifo.writeStart()
	var p uint32
	if p, err = e.WriteHeader(info); err != nil {
		return
	}
	e.offset = uint64(pos) + uint64(p)
	e.format = info.format
	e.depth = (info.bps + 7) / 8
	e.flenStd = (256 * (info.sps) / 245)
	e.flenLast = info.samples % e.flenStd
	e.frames = info.samples / e.flenStd
	if e.flenLast != 0 {
		e.frames++
	} else {
		e.flenLast = e.flenStd
	}
	e.rate = 0
	e.fifo.writeSkipBytes((e.frames + 1) * 4)
	e.seekTable = make([]uint64, e.frames)
	e.channels = int(info.nch)
	e.shiftBits = (4 - e.depth) << 3
	e.frameInit(0)
	return
}

func (e *Encoder) Close() {
	e.Finalize()
}

func (e *Encoder) Finalize() {
	e.fifo.writeDone()
	e.writeSeekTable()
}
