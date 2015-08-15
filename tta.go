package tta

import (
	"io"
)

type tta_info struct {
	format  uint32 // audio format
	nch     uint32 // number of channels
	bps     uint32 // bits per sample
	sps     uint32 // samplerate (sps)
	samples uint32 // data length in samples
}

type tta_filter interface {
	Decode(*int32)
	Encode(*int32)
}

type tta_filter_compat struct {
	index int32
	error int32
	round int32
	shift int32
	qm    [8]int32
	dx    [24]int32
	dl    [24]int32
}

type tta_filter_sse tta_filter_compat

type tta_adapt struct {
	k0   uint32
	k1   uint32
	sum0 uint32
	sum1 uint32
}

func (rice *tta_adapt) init(k0, k1 uint32) {
	rice.k0 = k0
	rice.k1 = k1
	rice.sum0 = shift_16[k0]
	rice.sum1 = shift_16[k1]
}

type tta_codec struct {
	filter tta_filter
	rice   tta_adapt
	prev   int32
}

type tta_fifo struct {
	buffer [TTA_FIFO_BUFFER_SIZE]byte
	pos    int32
	end    int32
	bcount uint32 // count of bits in cache
	bcache uint32 // bit cache
	crc    uint32
	count  uint32
	io     io.ReadWriteSeeker
}

type Decoder struct {
	codec        [MAX_NCH]tta_codec // 1 per channel
	channels     int                // number of channels/codecs
	data         [8]byte            // codec initialization data
	fifo         tta_fifo
	password_set bool     // password protection flag
	seek_allowed bool     // seek table flag
	seek_table   []uint64 // the playing position table
	format       uint32   // tta data format
	rate         uint32   // bitrate (kbps)
	offset       uint64   // data start position (header size, bytes)
	frames       uint32   // total count of frames
	depth        uint32   // bytes per sample
	flen_std     uint32   // default frame length in samples
	flen_last    uint32   // last frame length in samples
	flen         uint32   // current frame length in samples
	fnum         uint32   // currently playing frame index
	fpos         uint32   // the current position in frame
}

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

type Callback func(uint32, uint32, uint32)
