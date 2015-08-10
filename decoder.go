package tta

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
	return !this.fifo.read_crc32()
}

func (this *Decoder) set_password(pass string) {
	this.data = compute_key_digits([]byte(pass))
	this.password_set = true
}

func (this *Decoder) frame_init(frame uint32, seek_needed bool) error {
	if frame >= this.frames {
		return nil
	}
	shift := flt_set[this.depth-1]
	this.fnum = frame
	if seek_needed && this.seek_allowed {
		pos := this.seek_table[this.fnum]
		if pos != 0 && this.fifo.io.Seek(int64(pos)) < 0 {
			return TTA_SEEK_ERROR
		}
		this.fifo.read_start()
	}
	if this.fnum == this.frames-1 {
		this.flen = this.flen_last
	} else {
		this.flen = this.flen_std
	}
	for i := 0; i < this.decoder_len; i++ {
		this.decoder[i].fst.init(this.data, shift)
		this.decoder[i].rice.init(10, 10)
		this.decoder[i].prev = 0
	}
	this.fpos = 0
	this.fifo.reset()
	return nil
}

func (this *Decoder) frame_reset(frame uint32, iocb io_callback) {
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

func (this *Decoder) init_set_info(info *tta_info) error {
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

func (this *Decoder) init_get_info(info *tta_info, pos uint64) (err error) {
	if pos != 0 && this.fifo.io.Seek(int64(pos)) < 0 {
		err = TTA_SEEK_ERROR
		return
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
