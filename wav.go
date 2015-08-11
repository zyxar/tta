package tta

import (
	"errors"
	"os"
	"reflect"
	"unsafe"
)

const (
	RIFF_SIGN = (0x46464952)
	WAVE_SIGN = (0x45564157)
	fmt_SIGN  = (0x20746D66)
	data_SIGN = (0x61746164)

	WAVE_FORMAT_PCM        = 1
	WAVE_FORMAT_EXTENSIBLE = 0xFFFE
	PCM_BUFFER_LENGTH      = 5120
)

var PARTIAL_WRITTEN_ERROR = errors.New("partial written")

type WAVE_hdr struct {
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

type WAVE_subchunk_hdr struct {
	subchunk_id   uint32
	subchunk_size uint32
}

type WAVE_subformat struct {
	f1 uint32
	f2 uint16
	f3 uint16
	f4 [8]byte
}

type WAVE_ext_hdr struct {
	cb_size    uint16
	valid_bits uint16
	ch_mask    uint32
	est        WAVE_subformat
}

func (this *WAVE_hdr) toSlice() []byte {
	size := int(unsafe.Sizeof(*this))
	return *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(this)),
		Len:  int(size),
		Cap:  int(size),
	}))
}

func (this *WAVE_subchunk_hdr) toSlice() []byte {
	size := int(unsafe.Sizeof(*this))
	return *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(this)),
		Len:  int(size),
		Cap:  int(size),
	}))
}

func (this *WAVE_ext_hdr) toSlice() []byte {
	size := int(unsafe.Sizeof(*this))
	return *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(this)),
		Len:  int(size),
		Cap:  int(size),
	}))
}

func (this *WAVE_hdr) Read(infile *os.File) (subchunk_size uint32, err error) {
	var default_subchunk_size uint32 = 16
	b := this.toSlice()
	// Read WAVE header
	if _, err = infile.Read(b); err != nil {
		return
	}
	if this.audio_format == WAVE_FORMAT_EXTENSIBLE {
		wave_hdr_ex := WAVE_ext_hdr{}
		if _, err = infile.Read(wave_hdr_ex.toSlice()); err != nil {
			return
		}
		default_subchunk_size += uint32(unsafe.Sizeof(wave_hdr_ex))
		this.audio_format = uint16(wave_hdr_ex.est.f1)
	}

	// Skip extra format bytes
	if this.subchunk_size > default_subchunk_size {
		extra_len := this.subchunk_size - default_subchunk_size
		if _, err = infile.Seek(int64(extra_len), os.SEEK_SET); err != nil {
			return
		}
	}

	// Skip unsupported chunks
	subchunk_hdr := WAVE_subchunk_hdr{}
	for {
		if _, err = infile.Read(subchunk_hdr.toSlice()); err != nil {
			return
		}
		if subchunk_hdr.subchunk_id == data_SIGN {
			break
		}
		if _, err = infile.Seek(int64(subchunk_hdr.subchunk_size), os.SEEK_SET); err != nil {
			return
		}
	}
	subchunk_size = subchunk_hdr.subchunk_size
	return
}

func (this *WAVE_hdr) Write(fd *os.File, size uint32) (err error) {
	var write_len int
	// Write WAVE header
	if write_len, err = fd.Write(this.toSlice()); err != nil {
		return
	} else if write_len != int(unsafe.Sizeof(*this)) {
		err = PARTIAL_WRITTEN_ERROR
		return
	}
	// Write Subchunk header
	subchunk_hdr := WAVE_subchunk_hdr{data_SIGN, size}
	if write_len, err = fd.Write(subchunk_hdr.toSlice()); err != nil {
		return
	} else if write_len != int(unsafe.Sizeof(subchunk_hdr)) {
		err = PARTIAL_WRITTEN_ERROR
		return
	}
	return
}
