package tta

import (
	"io"
	"os"
)

var SSE_Enabled bool

func Decompress(infile, outfile *os.File, passwd string, cb Callback) (err error) {
	decoder := NewDecoder(infile)
	if len(passwd) > 0 {
		decoder.SetPassword(passwd)
	}
	info := tta_info{}
	if err = decoder.GetInfo(&info, 0); err != nil {
		return
	}
	smp_size := info.nch * ((info.bps + 7) / 8)
	data_size := info.samples * smp_size
	wave_hdr := WaveHeader{
		chunk_id:        RIFF_SIGN,
		chunk_size:      data_size + 36,
		format:          WAVE_SIGN,
		subchunk_id:     fmt_SIGN,
		subchunk_size:   16,
		audio_format:    1,
		num_channels:    uint16(info.nch),
		sample_rate:     info.sps,
		bits_per_sample: uint16(info.bps),
		byte_rate:       info.sps * smp_size,
		block_align:     uint16(smp_size),
	}
	if err = wave_hdr.Write(outfile, data_size); err != nil {
		return
	}
	buf_size := PCM_BUFFER_LENGTH * smp_size
	buffer := make([]byte, buf_size)
	var write_len int
	for {
		if write_len = int(uint32(decoder.ProcessStream(buffer, cb)) * smp_size); write_len == 0 {
			break
		}
		buf := buffer[:write_len]
		if write_len, err = outfile.Write(buf); err != nil {
			return
		} else if write_len != len(buf) {
			err = PARTIAL_WRITTEN_ERROR
			return
		}
	}
	return
}

func NewDecoder(iocb io.ReadWriteSeeker) *Decoder {
	dec := Decoder{}
	dec.fifo.io = iocb
	return &dec
}

func (this *Decoder) ProcessStream(out []byte, cb Callback) int32 {
	var cache [MAX_NCH]int32
	var value int32
	var ret int32 = 0
	i := 0
	out_ := out[:]
	for this.fpos < this.flen && len(out_) > 0 {
		value = this.fifo.get_value(&this.decoder[i].rice)
		// decompress stage 1: adaptive hybrid filter
		this.decoder[i].filter.Decode(&value)
		// decompress stage 2: fixed order 1 prediction
		value += ((this.decoder[i].prev * ((1 << 5) - 1)) >> 5)
		this.decoder[i].prev = value
		cache[i] = value
		if i < this.decoder_len-1 {
			i++
		} else {
			if this.decoder_len == 1 {
				write_buffer(value, out_, this.depth)
				out_ = out_[this.depth:]
			} else {
				k := i - 1
				cache[i] += cache[k] / 2
				for k > 0 {
					cache[k] = cache[k+1] - cache[k]
					k--
				}
				cache[k] = cache[k+1] - cache[k]
				for k <= i {
					write_buffer(cache[k], out_, this.depth)
					out_ = out_[this.depth:]
					k++
				}
			}
			i = 0
			this.fpos++
			ret++
		}
		if this.fpos == this.flen {
			// check frame crc
			crc_flag := !this.fifo.read_crc32()
			if crc_flag {
				for i := 0; i < len(out); i++ {
					out[i] = 0
				}
				if !this.seek_allowed {
					break
				}
			}
			this.fnum++

			// update dynamic info
			this.rate = (this.fifo.count << 3) / 1070
			if cb != nil {
				cb(this.rate, this.fnum, this.frames)
			}
			if this.fnum == this.frames {
				break
			}
			this.frame_init(this.fnum, crc_flag)
		}
	}
	return ret
}

func (this *Decoder) ProcessFrame(in_size uint32, out []byte) int32 {
	i := 0
	var cache [MAX_NCH]int32
	var value int32
	var ret int32 = 0
	out_ := out[:]
	for this.fifo.count < in_size && len(out_) > 0 {
		value = this.fifo.get_value(&this.decoder[i].rice)
		// decompress stage 1: adaptive hybrid filter
		this.decoder[i].filter.Decode(&value)
		// decompress stage 2: fixed order 1 prediction
		value += ((this.decoder[i].prev * ((1 << 5) - 1)) >> 5)
		this.decoder[i].prev = value
		cache[i] = value
		if i < this.decoder_len-1 {
			i++
		} else {
			if this.decoder_len == 1 {
				write_buffer(value, out_, this.depth)
				out_ = out_[this.depth:]
			} else {
				j := i
				k := i - 1
				cache[i] += cache[k] / 2
				for k > 0 {
					cache[k] = cache[j] - cache[k]
					j--
					k--
				}
				cache[k] = cache[j] - cache[k]
				for k <= i {
					write_buffer(cache[k], out_, this.depth)
					out_ = out_[this.depth:]
					k++
				}
			}
			i = 0
			this.fpos++
			ret++
		}

		if this.fpos == this.flen || this.fifo.count == in_size-4 {
			// check frame crc
			if !this.fifo.read_crc32() {
				for i := 0; i < len(out); i++ {
					out[i] = 0
				}
			}
			// update dynamic info
			this.rate = (this.fifo.count << 3) / 1070
			break
		}
	}
	return ret
}

func (this *Decoder) read_seek_table() bool {
	if this.seek_table == nil {
		return false
	}
	this.fifo.reset()
	tmp := this.offset + uint64(this.frames+1)*4
	for i := uint32(0); i < this.frames; i++ {
		this.seek_table[i] = tmp
		tmp += uint64(this.fifo.read_uint32())
	}
	return this.fifo.read_crc32()
}

func (this *Decoder) SetPassword(pass string) {
	this.data = compute_key_digits([]byte(pass))
	this.password_set = true
}

func (this *Decoder) frame_init(frame uint32, seek_needed bool) (err error) {
	if frame >= this.frames {
		return
	}
	shift := flt_set[this.depth-1]
	this.fnum = frame
	if seek_needed && this.seek_allowed {
		pos := this.seek_table[this.fnum]
		if pos != 0 {
			if _, err = this.fifo.io.Seek(int64(pos), os.SEEK_SET); err != nil {
				return TTA_SEEK_ERROR
			}
		}
		this.fifo.read_start()
	}
	if this.fnum == this.frames-1 {
		this.flen = this.flen_last
	} else {
		this.flen = this.flen_std
	}
	for i := 0; i < this.decoder_len; i++ {
		if SSE_Enabled {
			this.decoder[i].filter = NewSSEFilter(this.data, shift)
		} else {
			this.decoder[i].filter = NewCompatibleFilter(this.data, shift)
		}
		this.decoder[i].rice.init(10, 10)
		this.decoder[i].prev = 0
	}
	this.fpos = 0
	this.fifo.reset()
	return
}

func (this *Decoder) frame_reset(frame uint32, iocb io.ReadWriteSeeker) {
	this.fifo.io = iocb
	this.fifo.read_start()
	this.frame_init(frame, false)
}

func (this *Decoder) set_position(seconds uint32) (new_pos uint32, err error) {
	var frame uint32 = (245 * (seconds) / 256)
	new_pos = (256 * (frame) / 245)
	if !this.seek_allowed || frame >= this.frames {
		err = TTA_SEEK_ERROR
		return
	}
	this.frame_init(frame, true)
	return
}

func (this *Decoder) SetInfo(info *tta_info) error {
	if info.format > 2 ||
		info.bps < MIN_BPS ||
		info.bps > MAX_BPS ||
		info.nch > MAX_NCH {
		return TTA_FORMAT_ERROR
	}
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
	this.decoder_len = int(info.nch)
	this.fifo.read_start()
	this.frame_init(0, false)
	return nil
}

func (this *Decoder) GetInfo(info *tta_info, pos uint64) (err error) {
	if pos != 0 {
		if _, err = this.fifo.io.Seek(int64(pos), os.SEEK_SET); err != nil {
			err = TTA_SEEK_ERROR
			return
		}
	}
	this.fifo.read_start()
	var p uint32
	if p, err = this.fifo.read_tta_header(info); err != nil {
		return
	}
	pos += uint64(p)
	if info.format > 2 ||
		info.bps < MIN_BPS ||
		info.bps > MAX_BPS ||
		info.nch > MAX_NCH {
		return TTA_FORMAT_ERROR
	}
	if info.format == TTA_FORMAT_ENCRYPTED {
		if !this.password_set {
			return TTA_PASSWORD_ERROR
		}
	}
	this.offset = pos
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
	this.seek_table = make([]uint64, this.frames)
	this.seek_allowed = this.read_seek_table()
	this.decoder_len = int(info.nch)
	this.frame_init(0, false)
	return
}
