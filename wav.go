package tta

import (
	"os"
	"reflect"
	"unsafe"
)

const (
	_RIFF_SIGN = 0x46464952
	_WAVE_SIGN = 0x45564157
	_FMT_SIGN  = 0x20746D66
	_DATA_SIGN = 0x61746164

	_WAVE_FORMAT_PCM        = 1
	_WAVE_FORMAT_EXTENSIBLE = 0xFFFE
)

type WaveHeader struct {
	chunk_id        uint32
	chunk_size      uint32
	format          uint32
	subchunk_id     uint32
	subchunk_size   uint32
	audio_format    uint16
	num_channels    uint16
	sample_rate     uint32
	byte_rate       uint32
	block_align     uint16
	bits_per_sample uint16
}

type WaveSubchunkHeader struct {
	subchunk_id   uint32
	subchunk_size uint32
}

type WaveExtHeader struct {
	cb_size    uint16
	valid_bits uint16
	ch_mask    uint32
	est        struct {
		f1 uint32
		f2 uint16
		f3 uint16
		f4 [8]byte
	} // WaveSubformat
}

func (this *WaveHeader) Bytes() []byte {
	size := int(unsafe.Sizeof(*this))
	return *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(this)),
		Len:  size,
		Cap:  size,
	}))
}

func (this *WaveSubchunkHeader) Bytes() []byte {
	size := int(unsafe.Sizeof(*this))
	return *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(this)),
		Len:  size,
		Cap:  size,
	}))
}

func (this *WaveExtHeader) Bytes() []byte {
	size := int(unsafe.Sizeof(*this))
	return *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(this)),
		Len:  size,
		Cap:  size,
	}))
}

func (this *WaveHeader) Read(fd *os.File) (subchunk_size uint32, err error) {
	var default_subchunk_size uint32 = 16
	b := this.Bytes()
	var read_len int
	// Read WAVE header
	if read_len, err = fd.Read(b); err != nil {
		return
	} else if read_len != len(b) {
		err = PARTIAL_READ_ERROR
		return
	}
	if this.audio_format == _WAVE_FORMAT_EXTENSIBLE {
		wave_hdr_ex := WaveExtHeader{}
		if read_len, err = fd.Read(wave_hdr_ex.Bytes()); err != nil {
			return
		} else if read_len != int(unsafe.Sizeof(wave_hdr_ex)) {
			err = PARTIAL_READ_ERROR
			return
		}
		default_subchunk_size += uint32(unsafe.Sizeof(wave_hdr_ex))
		this.audio_format = uint16(wave_hdr_ex.est.f1)
	}

	// Skip extra format bytes
	if this.subchunk_size > default_subchunk_size {
		extra_len := this.subchunk_size - default_subchunk_size
		if _, err = fd.Seek(int64(extra_len), os.SEEK_SET); err != nil {
			return
		}
	}

	// Skip unsupported chunks
	subchunk_hdr := WaveSubchunkHeader{}
	for {
		if read_len, err = fd.Read(subchunk_hdr.Bytes()); err != nil {
			return
		} else if read_len != int(unsafe.Sizeof(subchunk_hdr)) {
			err = PARTIAL_READ_ERROR
			return
		}
		if subchunk_hdr.subchunk_id == _DATA_SIGN {
			break
		}
		if _, err = fd.Seek(int64(subchunk_hdr.subchunk_size), os.SEEK_SET); err != nil {
			return
		}
	}
	subchunk_size = subchunk_hdr.subchunk_size
	return
}

func (this *WaveHeader) Write(fd *os.File, size uint32) (err error) {
	var write_len int
	// Write WAVE header
	if write_len, err = fd.Write(this.Bytes()); err != nil {
		return
	} else if write_len != int(unsafe.Sizeof(*this)) {
		err = PARTIAL_WRITTEN_ERROR
		return
	}
	// Write Subchunk header
	subchunk_hdr := WaveSubchunkHeader{_DATA_SIGN, size}
	if write_len, err = fd.Write(subchunk_hdr.Bytes()); err != nil {
		return
	} else if write_len != int(unsafe.Sizeof(subchunk_hdr)) {
		err = PARTIAL_WRITTEN_ERROR
		return
	}
	return
}
