package wave

import (
	"io"
	"os"
	"reflect"
	"unsafe"
)

const (
	riff        = 0x46464952 // "RIFF"
	fmtChunkId  = 0x20746D66 // "fmt "
	dataChunkId = 0x61746164 // "data"
	MAGIC       = 0x45564157 // "WAVE"

	formatPCM = 1
	formatExt = 0xFFFE

	szHeader         = 36 // unsafe.Sizeof(Header{})
	szSubchunkHeader = 8  // unsafe.Sizeof(SubchunkHeader{})
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
		ChunkId:       riff,
		ChunkSize:     dataSize + szHeader,
		Format:        MAGIC,
		SubchunkId:    fmtChunkId,
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
	return (h.ChunkId == riff) &&
		(h.Format == MAGIC) &&
		(h.NumChannels != 0) &&
		(h.NumChannels <= maxNCH) &&
		(h.BitsPerSample != 0) &&
		(h.BitsPerSample <= maxBPS)
}

func (h *Header) Bytes() []byte {
	return slice(unsafe.Sizeof(*h), unsafe.Pointer(h))
}

func (s *SubchunkHeader) Bytes() []byte {
	return slice(unsafe.Sizeof(*s), unsafe.Pointer(s))
}

func (e *ExtHeader) Bytes() []byte {
	return slice(unsafe.Sizeof(*e), unsafe.Pointer(e))
}

func ReadHeader(fd io.ReadSeeker) (h *Header, subchunkSize uint32, err error) {
	var defaultSubchunkSize uint32 = 16 // PCM
	h = &Header{}
	b := h.Bytes()
	// Read WAVE header
	if _, err = fd.Read(b); err != nil {
		return
	}
	if h.AudioFormat == formatExt {
		extheader := ExtHeader{}
		if _, err = fd.Read(extheader.Bytes()); err != nil {
			return
		}
		defaultSubchunkSize += uint32(unsafe.Sizeof(extheader))
		h.AudioFormat = uint16(extheader.est.f1)
	}

	// Skip extra format bytes
	if h.SubchunkSize > defaultSubchunkSize {
		extraLen := h.SubchunkSize - defaultSubchunkSize
		if _, err = fd.Seek(int64(extraLen), os.SEEK_CUR); err != nil {
			return
		}
	}

	// Skip unsupported chunks
	subchunkHdr := SubchunkHeader{}
	for {
		if _, err = fd.Read(subchunkHdr.Bytes()); err != nil {
			return
		}
		if subchunkHdr.SubchunkId == dataChunkId {
			break
		}
		if _, err = fd.Seek(int64(subchunkHdr.SubchunkSize), os.SEEK_CUR); err != nil {
			return
		}
	}
	subchunkSize = subchunkHdr.SubchunkSize
	return
}

func (h *Header) WriteTo(w io.Writer) (n int64, err error) {
	// Write WAVE header
	writeLen, err := w.Write(h.Bytes())
	n += int64(writeLen)
	if err != nil {
		return
	}
	// Write Subchunk header
	subchunkHdr := SubchunkHeader{dataChunkId, h.ChunkSize - szHeader}
	writeLen, err = w.Write(subchunkHdr.Bytes())
	n += int64(writeLen)
	return
}

func slice(sz uintptr, addr unsafe.Pointer) []byte {
	size := int(sz)
	return *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(addr),
		Len:  size,
		Cap:  size,
	}))
}
