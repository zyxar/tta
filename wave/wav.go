package wave

import (
	"errors"
	"io"
	"os"
	"reflect"
	"unsafe"
)

const (
	magicRiff   = 0x46464952
	magicFormat = 0x20746D66
	magicChunk  = 0x61746164
	MAGIC       = 0x45564157

	formatPCM = 1
	formatExt = 0xFFFE
)

var (
	errPartialWritten = errors.New("partial written")
	errPartialRead    = errors.New("partial read")
)

type Header struct {
	ChunkId       uint32
	ChunkSize     uint32
	Format        uint32
	SubchunkId    uint32
	SubchunkSize  uint32
	AudioFormat   uint16
	NumChannels   uint16
	SampleRate    uint32
	ByteRate      uint32
	BlockAlign    uint16
	BitsPerSample uint16
}

type SubchunkHeader struct {
	SubchunkId   uint32
	SubchunkSize uint32
}

type ExtHeader struct {
	cbSize    uint16
	validBits uint16
	chMask    uint32
	est       struct {
		f1 uint32
		f2 uint16
		f3 uint16
		f4 [8]byte
	} // WaveSubformat
}

func NewHeader(dataSize uint32, nch uint16, sps uint32, bps uint16, smpSize uint16) *Header {
	return &Header{
		ChunkId:       magicRiff,
		ChunkSize:     dataSize + 36,
		Format:        MAGIC,
		SubchunkId:    magicFormat,
		SubchunkSize:  16,
		AudioFormat:   1,
		NumChannels:   nch,
		SampleRate:    sps,
		ByteRate:      sps * uint32(smpSize),
		BlockAlign:    smpSize,
		BitsPerSample: bps,
	}
}

func (h *Header) Validate(maxNCH, maxBPS uint16) bool {
	return (h.ChunkId == magicRiff) &&
		(h.Format == MAGIC) &&
		(h.NumChannels != 0) &&
		(h.NumChannels <= maxNCH) &&
		(h.BitsPerSample != 0) &&
		(h.BitsPerSample <= maxBPS)
}

func (w *Header) Bytes() []byte {
	size := int(unsafe.Sizeof(*w))
	return *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(w)),
		Len:  size,
		Cap:  size,
	}))
}

func (w *SubchunkHeader) Bytes() []byte {
	size := int(unsafe.Sizeof(*w))
	return *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(w)),
		Len:  size,
		Cap:  size,
	}))
}

func (w *ExtHeader) Bytes() []byte {
	size := int(unsafe.Sizeof(*w))
	return *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(w)),
		Len:  size,
		Cap:  size,
	}))
}

func (w *Header) Read(fd io.ReadSeeker) (subchunkSize uint32, err error) {
	var defaultSubchunkSize uint32 = 16
	b := w.Bytes()
	var readLen int
	// Read WAVE header
	if readLen, err = fd.Read(b); err != nil {
		return
	} else if readLen != len(b) {
		err = errPartialRead
		return
	}
	if w.AudioFormat == formatExt {
		waveHdrEx := ExtHeader{}
		if readLen, err = fd.Read(waveHdrEx.Bytes()); err != nil {
			return
		} else if readLen != int(unsafe.Sizeof(waveHdrEx)) {
			err = errPartialRead
			return
		}
		defaultSubchunkSize += uint32(unsafe.Sizeof(waveHdrEx))
		w.AudioFormat = uint16(waveHdrEx.est.f1)
	}

	// Skip extra format bytes
	if w.SubchunkSize > defaultSubchunkSize {
		extraLen := w.SubchunkSize - defaultSubchunkSize
		if _, err = fd.Seek(int64(extraLen), os.SEEK_SET); err != nil {
			return
		}
	}

	// Skip unsupported chunks
	subchunkHdr := SubchunkHeader{}
	for {
		if readLen, err = fd.Read(subchunkHdr.Bytes()); err != nil {
			return
		} else if readLen != int(unsafe.Sizeof(subchunkHdr)) {
			err = errPartialRead
			return
		}
		if subchunkHdr.SubchunkId == magicChunk {
			break
		}
		if _, err = fd.Seek(int64(subchunkHdr.SubchunkSize), os.SEEK_SET); err != nil {
			return
		}
	}
	subchunkSize = subchunkHdr.SubchunkSize
	return
}

func (w *Header) Write(fd io.Writer) (err error) {
	var writeLen int
	// Write WAVE header
	if writeLen, err = fd.Write(w.Bytes()); err != nil {
		return
	} else if writeLen != int(unsafe.Sizeof(*w)) {
		err = errPartialWritten
		return
	}
	// Write Subchunk header
	subchunkHdr := SubchunkHeader{magicChunk, w.ChunkSize - 36}
	if writeLen, err = fd.Write(subchunkHdr.Bytes()); err != nil {
		return
	} else if writeLen != int(unsafe.Sizeof(subchunkHdr)) {
		err = errPartialWritten
		return
	}
	return
}
