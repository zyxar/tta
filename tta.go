package tta

type tta_info struct {
	format  uint32 // audio format
	nch     uint32 // number of channels
	bps     uint32 // bits per sample
	sps     uint32 // samplerate (sps)
	samples uint32 // data length in samples
}

type tta_fltst struct {
	index int32
	error int32
	round int32
	shift int32
	qm    [8]int32
	dx    [24]int32
	dl    [24]int32
}

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
	fst  tta_fltst
	rice tta_adapt
	prev int32
}

type tta_fifo struct {
	buffer [TTA_FIFO_BUFFER_SIZE]byte
	pos    int32
	end    int32
	bcount uint32 // count of bits in cache
	bcache uint32 // bit cache
	crc    uint32
	count  uint32
	io     io_callback
}

type Decoder struct {
	decoder      [MAX_NCH]tta_codec // decoder (1 per channel)
	data         [8]int8            // decoder initialization data
	fifo         tta_fifo
	decoder_last *tta_codec
	password_set bool    // password protection flag
	seek_table   *uint64 // the playing position table
	format       uint32  // tta data format
	rate         uint32  // bitrate (kbps)
	offset       uint64  // data start position (header size, bytes)
	frames       uint32  // total count of frames
	depth        uint32  // bytes per sample
	flen_std     uint32  // default frame length in samples
	flen_last    uint32  // last frame length in samples
	flen         uint32  // current frame length in samples
	fnum         uint32  // currently playing frame index
	fpos         uint32  // the current position in frame
}

type Encoder struct {
	encoder      [MAX_NCH]tta_codec // encoder (1 per channel)
	data         [8]int8            // encoder initialization data
	fifo         tta_fifo
	encoder_last *tta_codec
	seek_table   *uint64 // the playing position table
	format       uint32  // tta data format
	rate         uint32  // bitrate (kbps)
	offset       uint64  // data start position (header size, bytes)
	frames       uint32  // total count of frames
	depth        uint32  // bytes per sample
	flen_std     uint32  // default frame length in samples
	flen_last    uint32  // last frame length in samples
	flen         uint32  // current frame length in samples
	fnum         uint32  // currently playing frame index
	fpos         uint32  // the current position in frame
	shift_bits   uint32  // packing int to pcm
}

type Callback func(uint32, uint32, uint32)
type io_callback interface {
	Read(p []byte) int32
	Write(p []byte) int32
	Seek(offset int64) int64
}
