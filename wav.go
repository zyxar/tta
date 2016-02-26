package tta

import (
	"io"
	"os"
	"reflect"
	"unsafe"
)

const (
	riffSign = 0x46464952
	waveSign = 0x45564157
	fmtSign  = 0x20746D66
	dataSign = 0x61746164

	waveFormatPcm        = 1
	waveFormatExtensible = 0xFFFE
)

type WaveHeader struct {
	chunkId       uint32
	chunkSize     uint32
	format        uint32
	subchunkId    uint32
	subchunkSize  uint32
	audioFormat   uint16
	numChannels   uint16
	sampleRate    uint32
	byteRate      uint32
	blockAlign    uint16
	bitsPerSample uint16
}

type WaveSubchunkHeader struct {
	subchunkId   uint32
	subchunkSize uint32
}

type WaveExtHeader struct {
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

func (w *WaveHeader) Bytes() []byte {
	size := int(unsafe.Sizeof(*w))
	return *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(w)),
		Len:  size,
		Cap:  size,
	}))
}

func (w *WaveSubchunkHeader) Bytes() []byte {
	size := int(unsafe.Sizeof(*w))
	return *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(w)),
		Len:  size,
		Cap:  size,
	}))
}

func (w *WaveExtHeader) Bytes() []byte {
	size := int(unsafe.Sizeof(*w))
	return *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(w)),
		Len:  size,
		Cap:  size,
	}))
}

func (w *WaveHeader) Read(fd io.ReadSeeker) (subchunkSize uint32, err error) {
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
	if w.audioFormat == waveFormatExtensible {
		waveHdrEx := WaveExtHeader{}
		if readLen, err = fd.Read(waveHdrEx.Bytes()); err != nil {
			return
		} else if readLen != int(unsafe.Sizeof(waveHdrEx)) {
			err = errPartialRead
			return
		}
		defaultSubchunkSize += uint32(unsafe.Sizeof(waveHdrEx))
		w.audioFormat = uint16(waveHdrEx.est.f1)
	}

	// Skip extra format bytes
	if w.subchunkSize > defaultSubchunkSize {
		extraLen := w.subchunkSize - defaultSubchunkSize
		if _, err = fd.Seek(int64(extraLen), os.SEEK_SET); err != nil {
			return
		}
	}

	// Skip unsupported chunks
	subchunkHdr := WaveSubchunkHeader{}
	for {
		if readLen, err = fd.Read(subchunkHdr.Bytes()); err != nil {
			return
		} else if readLen != int(unsafe.Sizeof(subchunkHdr)) {
			err = errPartialRead
			return
		}
		if subchunkHdr.subchunkId == dataSign {
			break
		}
		if _, err = fd.Seek(int64(subchunkHdr.subchunkSize), os.SEEK_SET); err != nil {
			return
		}
	}
	subchunkSize = subchunkHdr.subchunkSize
	return
}

func (w *WaveHeader) Write(fd io.Writer, size uint32) (err error) {
	var writeLen int
	// Write WAVE header
	if writeLen, err = fd.Write(w.Bytes()); err != nil {
		return
	} else if writeLen != int(unsafe.Sizeof(*w)) {
		err = errPartialWritten
		return
	}
	// Write Subchunk header
	subchunkHdr := WaveSubchunkHeader{dataSign, size}
	if writeLen, err = fd.Write(subchunkHdr.Bytes()); err != nil {
		return
	} else if writeLen != int(unsafe.Sizeof(subchunkHdr)) {
		err = errPartialWritten
		return
	}
	return
}
