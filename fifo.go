package tta

func (s *tta_fifo) read_byte() (v byte) {
	if s.pos >= s.end {
		v, _ := s.io.Read(s.buffer[:]) // FIXME: handle this error
		s.end = int32(v)
		s.pos = 0
	}
	s.crc = crc32_table[(s.crc^uint32(s.buffer[s.pos]))&0xFF] ^ (s.crc >> 8)
	s.count++
	v = s.buffer[s.pos]
	s.pos++
	return
}

func (s *tta_fifo) read_uint16() (v uint16) {
	v = 0
	v |= uint16(s.read_byte())
	v |= uint16(s.read_byte()) << 8
	return
}

func (s *tta_fifo) read_uint32() (v uint32) {
	v = 0
	v |= uint32(s.read_byte())
	v |= uint32(s.read_byte()) << 8
	v |= uint32(s.read_byte()) << 16
	v |= uint32(s.read_byte()) << 24
	return
}

func (s *tta_fifo) read_crc32() bool {
	crc := s.crc ^ 0xFFFFFFFF
	return crc == s.read_uint32()
}

func (s *tta_fifo) read_start() {
	s.pos = s.end
}

func (s *tta_fifo) reset() {
	s.crc = 0xFFFFFFFF
	s.bcache = 0
	s.bcount = 0
	s.count = 0
}

func (s *tta_fifo) read_skip_bytes(size uint32) {
	for size > 0 {
		size--
		s.read_byte()
	}
}

func (s *tta_fifo) skip_id3v2() (size uint32) {
	s.reset()
	if 'I' != s.read_byte() || 'D' != s.read_byte() || '3' != s.read_byte() {
		s.pos = 0
		return 0
	}
	s.pos += 2 // skip version bytes
	if s.read_byte()&0x10 != 0 {
		size += 10
	}
	size += uint32(s.read_byte() & 0x7F)
	size = (size << 7) | uint32(s.read_byte()&0x7F)
	size = (size << 7) | uint32(s.read_byte()&0x7F)
	size = (size << 7) | uint32(s.read_byte()&0x7F)
	s.read_skip_bytes(size)
	size += 10
	return
}

func (s *tta_fifo) read_tta_header(info *tta_info) (uint32, error) {
	size := s.skip_id3v2()
	s.reset()
	if 'T' != s.read_byte() ||
		'T' != s.read_byte() ||
		'A' != s.read_byte() ||
		'1' != s.read_byte() {
		return 0, TTA_FORMAT_ERROR
	}
	info.format = uint32(s.read_uint16())
	info.nch = uint32(s.read_uint16())
	info.bps = uint32(s.read_uint16())
	info.sps = s.read_uint32()
	info.samples = s.read_uint32()
	if !s.read_crc32() {
		return 0, TTA_FILE_ERROR
	}
	size += 22
	return size, nil
}

func (s *tta_fifo) get_value(rice *tta_adapt) (value int32) {
	if s.bcache^bit_mask[s.bcount] == 0 {
		value += int32(s.bcount)
		s.bcache = uint32(s.read_byte())
		s.bcount = 8
		for s.bcache == 0xFF {
			value += 8
			s.bcache = uint32(s.read_byte())
		}
	}

	for (s.bcache & 1) != 0 {
		value++
		s.bcache >>= 1
		s.bcount--
	}
	s.bcache >>= 1
	s.bcount--

	var level, k, tmp uint32
	if value != 0 {
		level = 1
		k = rice.k1
		value--
	} else {
		level = 0
		k = rice.k0
	}
	if k != 0 {
		for s.bcount < k {
			tmp = uint32(s.read_byte())
			s.bcache |= tmp << s.bcount
			s.bcount += 8
		}
		value = (value << k) + int32(s.bcache&bit_mask[k])
		s.bcache >>= k
		s.bcount -= k
		s.bcache &= bit_mask[s.bcount]
	}
	if level != 0 {
		rice.sum1 += uint32(value) - (rice.sum1 >> 4)
		if rice.k1 > 0 && rice.sum1 < shift_16[rice.k1] {
			rice.k1--
		} else if rice.sum1 > shift_16[rice.k1+1] {
			rice.k1++
		}
		value += int32(bit_shift[rice.k0])
	}

	rice.sum0 += uint32(value) - (rice.sum0 >> 4)
	if rice.k0 > 0 && rice.sum0 < shift_16[rice.k0] {
		rice.k0--
	} else if rice.sum0 > shift_16[rice.k0+1] {
		rice.k0++
	}
	// ((x & 1)?((x + 1) >> 1):(-x >> 1))
	if value&1 != 0 {
		value = (value + 1) >> 1
	} else {
		value = -value >> 1
	}
	return
}

func (s *tta_fifo) write_start() {
	s.pos = 0
}

func (s *tta_fifo) write_done() error {
	if s.pos > 0 {
		if n, err := s.io.Write(s.buffer[:s.pos]); err != nil || n != int(s.pos) {
			return TTA_WRITE_ERROR
		}
		s.pos = 0
	}
	return nil
}

func (s *tta_fifo) write_byte(v byte) error {
	if s.pos == TTA_FIFO_BUFFER_SIZE {
		if n, err := s.io.Write(s.buffer[:]); err != nil || n != TTA_FIFO_BUFFER_SIZE {
			return TTA_WRITE_ERROR
		}
		s.pos = 0
	}
	s.crc = crc32_table[(s.crc^uint32(v))&0xFF] ^ (s.crc >> 8)
	s.count++
	s.buffer[s.pos] = v
	s.pos++
	return nil
}

func (s *tta_fifo) write_uint16(v uint16) error {
	if err := s.write_byte(byte(v)); err != nil {
		return err
	}
	if err := s.write_byte(byte(v >> 8)); err != nil {
		return err
	}
	return nil
}

func (s *tta_fifo) write_uint32(v uint32) error {
	if err := s.write_byte(byte(v)); err != nil {
		return err
	}
	if err := s.write_byte(byte(v >> 8)); err != nil {
		return err
	}
	if err := s.write_byte(byte(v >> 16)); err != nil {
		return err
	}
	return s.write_byte(byte(v >> 24))
}

func (s *tta_fifo) write_crc32() error {
	return s.write_uint32(s.crc ^ 0xFFFFFFFF)
}

func (s *tta_fifo) write_skip_bytes(size uint32) error {
	for size > 0 {
		if err := s.write_byte(0); err != nil {
			return err
		}
		size--
	}
	return nil
}

func (s *tta_fifo) write_tta_header(info *tta_info) (size uint32, err error) {
	s.reset()

	// write TTA1 signature
	if err = s.write_byte('T'); err != nil {
		return
	}
	if err = s.write_byte('T'); err != nil {
		return
	}
	if err = s.write_byte('A'); err != nil {
		return
	}
	if err = s.write_byte('1'); err != nil {
		return
	}

	if err = s.write_uint16(uint16(info.format)); err != nil {
		return
	}
	if err = s.write_uint16(uint16(info.nch)); err != nil {
		return
	}
	if err = s.write_uint16(uint16(info.bps)); err != nil {
		return
	}
	if err = s.write_uint32(info.sps); err != nil {
		return
	}
	if err = s.write_uint32(info.samples); err != nil {
		return
	}

	if err = s.write_crc32(); err != nil {
		return
	}
	size = 22
	return
}

func (s *tta_fifo) put_value(rice *tta_adapt, value int32) {
	var k, unary, outval uint32
	if value > 0 {
		outval = (uint32(value) << 1) - 1
	} else {
		outval = uint32(-value) << 1
	}
	// encode Rice unsigned
	k = rice.k0
	rice.sum0 += outval - (rice.sum0 >> 4)
	if rice.k0 > 0 && rice.sum0 < shift_16[rice.k0] {
		rice.k0--
	} else if rice.sum0 > shift_16[rice.k0+1] {
		rice.k0++
	}

	if outval >= bit_shift[k] {
		outval -= bit_shift[k]
		k = rice.k1
		rice.sum1 += outval - (rice.sum1 >> 4)
		if rice.k1 > 0 && rice.sum1 < shift_16[rice.k1] {
			rice.k1--
		} else if rice.sum1 > shift_16[rice.k1+1] {
			rice.k1++
		}
		unary = 1 + (outval >> k)
	} else {
		unary = 0
	}

	for { // put unary
		for s.bcount >= 8 {
			s.write_byte(byte(s.bcache))
			s.bcache >>= 8
			s.bcount -= 8
		}
		if unary > 23 {
			s.bcache |= bit_mask[23] << s.bcount
			s.bcount += 23
			unary -= 23
		} else {
			s.bcache |= bit_mask[unary] << s.bcount
			s.bcount += unary + 1
			unary = 0
		}
		if unary == 0 {
			break
		}
	}
	for s.bcount >= 8 { // put binary
		s.write_byte(byte(s.bcache))
		s.bcache >>= 8
		s.bcount -= 8
	}
	if k != 0 {
		s.bcache |= (outval & bit_mask[k]) << s.bcount
		s.bcount += k
	}
}

func (s *tta_fifo) flush_bit_cache() {
	for s.bcount > 0 {
		s.write_byte(byte(s.bcache))
		s.bcache >>= 8
		if s.bcount > 8 {
			s.bcount -= 8
		} else {
			break
		}
	}
	s.write_crc32()
}
