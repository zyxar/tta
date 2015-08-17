package tta

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
