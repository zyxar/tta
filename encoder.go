package tta

import (
	"fmt"
	"io"
	"os"
)

type Encoder struct {
	codec      [MAX_NCH]tta_codec // 1 per channel
	channels   int                // number of channels/codecs
	data       [8]byte            // codec initialization data
	fifo       tta_fifo
	seek_table []uint64 // the playing position table
	format     uint32   // tta data format
	rate       uint32   // bitrate (kbps)
	offset     uint64   // data start position (header size, bytes)
	frames     uint32   // total count of frames
	depth      uint32   // bytes per sample
	flen_std   uint32   // default frame length in samples
	flen_last  uint32   // last frame length in samples
	flen       uint32   // current frame length in samples
	fnum       uint32   // currently playing frame index
	fpos       uint32   // the current position in frame
	shift_bits uint32   // packing int to pcm
}

func Compress(infile, outfile io.ReadWriteSeeker, passwd string, cb Callback) (err error) {
	wave_hdr := WaveHeader{}
	var data_size uint32
	if data_size, err = wave_hdr.Read(infile); err != nil {
		err = TTA_READ_ERROR
		return
	} else if data_size >= 0x7FFFFFFF {
		err = fmt.Errorf("incorrect data size info in wav file: %x", data_size)
		return
	}
	if (wave_hdr.chunk_id != _RIFF_SIGN) ||
		(wave_hdr.format != _WAVE_SIGN) ||
		(wave_hdr.num_channels == 0) ||
		(wave_hdr.num_channels > MAX_NCH) ||
		(wave_hdr.bits_per_sample == 0) ||
		(wave_hdr.bits_per_sample > MAX_BPS) {
		err = TTA_FORMAT_ERROR
		return
	}
	encoder := NewEncoder(outfile)
	smp_size := uint32(wave_hdr.num_channels * ((wave_hdr.bits_per_sample + 7) / 8))
	info := tta_info{
		nch:     uint32(wave_hdr.num_channels),
		bps:     uint32(wave_hdr.bits_per_sample),
		sps:     wave_hdr.sample_rate,
		format:  TTA_FORMAT_SIMPLE,
		samples: data_size / smp_size,
	}
	if len(passwd) > 0 {
		encoder.SetPassword(passwd)
		info.format = TTA_FORMAT_ENCRYPTED
	}
	buf_size := PCM_BUFFER_LENGTH * smp_size
	buffer := make([]byte, buf_size)
	if err = encoder.SetInfo(&info, 0); err != nil {
		return
	}
	var read_len int = 0
	for data_size > 0 {
		if buf_size >= data_size {
			buf_size = data_size
		}
		if read_len, err = infile.Read(buffer[:buf_size]); err != nil || read_len != int(buf_size) {
			err = TTA_READ_ERROR
			return
		}
		encoder.ProcessStream(buffer[:buf_size], cb)
		data_size -= buf_size
	}
	encoder.Close()
	return
}

func NewEncoder(iocb io.ReadWriteSeeker) *Encoder {
	enc := Encoder{}
	enc.fifo.io = iocb
	return &enc
}

func (this *Encoder) ProcessStream(in []byte, cb Callback) {
	if len(in) == 0 {
		return
	}
	var res, curr, next, tmp int32
	next = read_buffer(in, this.depth)
	in = in[this.depth:]
	tmp = next << this.shift_bits
	i := 0
	index := 0
	for {
		curr = next
		if index < len(in) {
			next = read_buffer(in[index:], this.depth)
			tmp = next << this.shift_bits
		} else {
			next = 0
			tmp = 0
		}
		index += int(this.depth)
		// transform data
		if this.channels > 1 {
			if i < this.channels-1 {
				res = next - curr
				curr = res
			} else {
				curr -= res / 2
			}
		}
		// compress stage 1: fixed order 1 prediction
		tmp = curr
		curr -= ((this.codec[i].prev * ((1 << 5) - 1)) >> 5)
		this.codec[i].prev = tmp
		// compress stage 2: adaptive hybrid filter
		this.codec[i].filter.Encode(&curr)
		this.fifo.put_value(&this.codec[i].rice, curr)
		if i < this.channels-1 {
			i++
		} else {
			i = 0
			this.fpos++
		}
		if this.fpos == this.flen {
			this.fifo.flush_bit_cache()
			this.seek_table[this.fnum] = uint64(this.fifo.count)
			this.fnum++
			// update dynamic info
			this.rate = (this.fifo.count << 3) / 1070
			if cb != nil {
				cb(this.rate, this.fnum, this.frames)
			}
			this.frame_init(this.fnum)
		}
		if index >= int(this.depth)+len(in) {
			break
		}
	}
}

func (this *Encoder) ProcessFrame(in []byte) {
	if len(in) == 0 {
		return
	}
	var res, curr, next, tmp int32
	next = read_buffer(in, this.depth)
	in = in[this.depth:]
	tmp = next << this.shift_bits
	i := 0
	index := 0
	for {
		curr = next
		if index < len(in) {
			next = read_buffer(in[index:], this.depth)
			tmp = next << this.shift_bits
		} else {
			next = 0
			tmp = 0
		}
		index += int(this.depth)
		// transform data
		if this.channels > 1 {
			if i < this.channels-1 {
				res = next - curr
				curr = res
			} else {
				curr -= res / 2
			}
		}
		// compress stage 1: fixed order 1 prediction
		tmp = curr
		curr -= ((this.codec[i].prev * ((1 << 5) - 1)) >> 5)
		this.codec[i].prev = tmp
		// compress stage 2: adaptive hybrid filter
		this.codec[i].filter.Encode(&curr)
		this.fifo.put_value(&this.codec[i].rice, curr)
		if i < this.channels-1 {
			i++
		} else {
			i = 0
			this.fpos++
		}
		if this.fpos == this.flen {
			this.fifo.flush_bit_cache()
			// update dynamic info
			this.rate = (this.fifo.count << 3) / 1070
			break
		}
		if index >= int(this.depth)+len(in) {
			break
		}
	}
}

func (this *Encoder) write_seek_table() (err error) {
	if this.seek_table == nil {
		return
	}
	if _, err = this.fifo.io.Seek(int64(this.offset), os.SEEK_SET); err != nil {
		return
	}
	this.fifo.write_start()
	this.fifo.reset()
	for i := uint32(0); i < this.frames; i++ {
		this.fifo.write_uint32(uint32(this.seek_table[i] & 0xFFFFFFFF))
	}
	this.fifo.write_crc32()
	this.fifo.write_done()
	return
}

func (this *Encoder) SetPassword(pass string) {
	this.data = compute_key_digits(convert_password(pass))
}

func (this *Encoder) frame_init(frame uint32) (err error) {
	if frame >= this.frames {
		return
	}
	shift := flt_set[this.depth-1]
	this.fnum = frame
	if this.fnum == this.frames-1 {
		this.flen = this.flen_last
	} else {
		this.flen = this.flen_std
	}
	// init entropy encoder
	for i := 0; i < this.channels; i++ {
		if SSE_Enabled {
			this.codec[i].filter = NewSSEFilter(this.data, shift)
		} else {
			this.codec[i].filter = NewCompatibleFilter(this.data, shift)
		}
		this.codec[i].rice.init(10, 10)
		this.codec[i].prev = 0
	}
	this.fpos = 0
	this.fifo.reset()
	return
}

func (this *Encoder) frame_reset(frame uint32, iocb io.ReadWriteSeeker) {
	this.fifo.io = iocb
	this.fifo.read_start()
	this.frame_init(frame)
}

func (this *Encoder) WriteHeader(info *tta_info) (size uint32, err error) {
	this.fifo.reset()
	// write TTA1 signature
	if err = this.fifo.write_byte('T'); err != nil {
		return
	}
	if err = this.fifo.write_byte('T'); err != nil {
		return
	}
	if err = this.fifo.write_byte('A'); err != nil {
		return
	}
	if err = this.fifo.write_byte('1'); err != nil {
		return
	}
	if err = this.fifo.write_uint16(uint16(info.format)); err != nil {
		return
	}
	if err = this.fifo.write_uint16(uint16(info.nch)); err != nil {
		return
	}
	if err = this.fifo.write_uint16(uint16(info.bps)); err != nil {
		return
	}
	if err = this.fifo.write_uint32(info.sps); err != nil {
		return
	}
	if err = this.fifo.write_uint32(info.samples); err != nil {
		return
	}
	if err = this.fifo.write_crc32(); err != nil {
		return
	}
	size = 22
	return

}

func (this *Encoder) SetInfo(info *tta_info, pos int64) (err error) {
	if info.format > 2 ||
		info.bps < MIN_BPS ||
		info.bps > MAX_BPS ||
		info.nch > MAX_NCH {
		return TTA_FORMAT_ERROR
	}
	// set start position if required
	if pos != 0 {
		if _, err = this.fifo.io.Seek(int64(pos), os.SEEK_SET); err != nil {
			err = TTA_SEEK_ERROR
			return
		}
	}
	this.fifo.write_start()
	var p uint32
	if p, err = this.WriteHeader(info); err != nil {
		return
	}
	this.offset = uint64(pos) + uint64(p)
	this.format = info.format
	this.depth = (info.bps + 7) / 8
	this.flen_std = (256 * (info.sps) / 245)
	this.flen_last = info.samples % this.flen_std
	this.frames = info.samples / this.flen_std
	if this.flen_last != 0 {
		this.frames += 1
	} else {
		this.flen_last = this.flen_std
	}
	this.rate = 0
	this.fifo.write_skip_bytes((this.frames + 1) * 4)
	this.seek_table = make([]uint64, this.frames)
	this.channels = int(info.nch)
	this.shift_bits = (4 - this.depth) << 3
	this.frame_init(0)
	return
}

func (this *Encoder) Close() {
	this.Finalize()
}

func (this *Encoder) Finalize() {
	this.fifo.write_done()
	this.write_seek_table()
}
