package tta

import (
	"io"
)

type fifo struct {
	buffer [fifoBufferSize]byte
	pos    int32
	end    int32
	bcount uint32 // count of bits in cache
	bcache uint32 // bit cache
	crc    uint32
	count  uint32
	io     io.ReadWriteSeeker
}

func (s *fifo) readByte() (v byte) {
	if s.pos >= s.end {
		v, _ := s.io.Read(s.buffer[:]) // FIXME: handle this error
		s.end = int32(v)
		s.pos = 0
	}
	s.crc = crc32Table[(s.crc^uint32(s.buffer[s.pos]))&0xFF] ^ (s.crc >> 8)
	s.count++
	v = s.buffer[s.pos]
	s.pos++
	return
}

func (s *fifo) readUint16() (v uint16) {
	v = 0
	v |= uint16(s.readByte())
	v |= uint16(s.readByte()) << 8
	return
}

func (s *fifo) readUint32() (v uint32) {
	v = 0
	v |= uint32(s.readByte())
	v |= uint32(s.readByte()) << 8
	v |= uint32(s.readByte()) << 16
	v |= uint32(s.readByte()) << 24
	return
}

func (s *fifo) readCrc32() bool {
	crc := s.crc ^ 0xFFFFFFFF
	return crc == s.readUint32()
}

func (s *fifo) readStart() {
	s.pos = s.end
}

func (s *fifo) reset() {
	s.crc = 0xFFFFFFFF
	s.bcache = 0
	s.bcount = 0
	s.count = 0
}

func (s *fifo) readSkipBytes(size uint32) {
	for size > 0 {
		size--
		s.readByte()
	}
}

func (s *fifo) skipId3v2() (size uint32) {
	s.reset()
	if 'I' != s.readByte() || 'D' != s.readByte() || '3' != s.readByte() {
		s.pos = 0
		return 0
	}
	s.pos += 2 // skip version bytes
	if s.readByte()&0x10 != 0 {
		size += 10
	}
	size += uint32(s.readByte() & 0x7F)
	size = (size << 7) | uint32(s.readByte()&0x7F)
	size = (size << 7) | uint32(s.readByte()&0x7F)
	size = (size << 7) | uint32(s.readByte()&0x7F)
	s.readSkipBytes(size)
	size += 10
	return
}

func (s *fifo) getValue(ad *adapter) (value int32) {
	if s.bcache^bitMask[s.bcount] == 0 {
		value += int32(s.bcount)
		s.bcache = uint32(s.readByte())
		s.bcount = 8
		for s.bcache == 0xFF {
			value += 8
			s.bcache = uint32(s.readByte())
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
		k = ad.k1
		value--
	} else {
		level = 0
		k = ad.k0
	}
	if k != 0 {
		for s.bcount < k {
			tmp = uint32(s.readByte())
			s.bcache |= tmp << s.bcount
			s.bcount += 8
		}
		value = (value << k) + int32(s.bcache&bitMask[k])
		s.bcache >>= k
		s.bcount -= k
		s.bcache &= bitMask[s.bcount]
	}
	if level != 0 {
		ad.sum1 += uint32(value) - (ad.sum1 >> 4)
		if ad.k1 > 0 && ad.sum1 < shift16[ad.k1] {
			ad.k1--
		} else if ad.sum1 > shift16[ad.k1+1] {
			ad.k1++
		}
		value += int32(bitShift[ad.k0])
	}

	ad.sum0 += uint32(value) - (ad.sum0 >> 4)
	if ad.k0 > 0 && ad.sum0 < shift16[ad.k0] {
		ad.k0--
	} else if ad.sum0 > shift16[ad.k0+1] {
		ad.k0++
	}
	// ((x & 1)?((x + 1) >> 1):(-x >> 1))
	if value&1 != 0 {
		value = (value + 1) >> 1
	} else {
		value = -value >> 1
	}
	return
}

func (s *fifo) writeStart() {
	s.pos = 0
}

func (s *fifo) writeDone() error {
	if s.pos > 0 {
		if n, err := s.io.Write(s.buffer[:s.pos]); err != nil || n != int(s.pos) {
			return errWrite
		}
		s.pos = 0
	}
	return nil
}

func (s *fifo) writeByte(v byte) error {
	if s.pos == fifoBufferSize {
		if n, err := s.io.Write(s.buffer[:]); err != nil || n != fifoBufferSize {
			return errWrite
		}
		s.pos = 0
	}
	s.crc = crc32Table[(s.crc^uint32(v))&0xFF] ^ (s.crc >> 8)
	s.count++
	s.buffer[s.pos] = v
	s.pos++
	return nil
}

func (s *fifo) writeUint16(v uint16) error {
	if err := s.writeByte(byte(v)); err != nil {
		return err
	}
	if err := s.writeByte(byte(v >> 8)); err != nil {
		return err
	}
	return nil
}

func (s *fifo) writeUint32(v uint32) error {
	if err := s.writeByte(byte(v)); err != nil {
		return err
	}
	if err := s.writeByte(byte(v >> 8)); err != nil {
		return err
	}
	if err := s.writeByte(byte(v >> 16)); err != nil {
		return err
	}
	return s.writeByte(byte(v >> 24))
}

func (s *fifo) writeCrc32() error {
	return s.writeUint32(s.crc ^ 0xFFFFFFFF)
}

func (s *fifo) writeSkipBytes(size uint32) error {
	for size > 0 {
		if err := s.writeByte(0); err != nil {
			return err
		}
		size--
	}
	return nil
}

func (s *fifo) putValue(ad *adapter, value int32) {
	var k, unary, outval uint32
	if value > 0 {
		outval = (uint32(value) << 1) - 1
	} else {
		outval = uint32(-value) << 1
	}
	// encode Rice unsigned
	k = ad.k0
	ad.sum0 += outval - (ad.sum0 >> 4)
	if ad.k0 > 0 && ad.sum0 < shift16[ad.k0] {
		ad.k0--
	} else if ad.sum0 > shift16[ad.k0+1] {
		ad.k0++
	}

	if outval >= bitShift[k] {
		outval -= bitShift[k]
		k = ad.k1
		ad.sum1 += outval - (ad.sum1 >> 4)
		if ad.k1 > 0 && ad.sum1 < shift16[ad.k1] {
			ad.k1--
		} else if ad.sum1 > shift16[ad.k1+1] {
			ad.k1++
		}
		unary = 1 + (outval >> k)
	} else {
		unary = 0
	}

	for { // put unary
		for s.bcount >= 8 {
			s.writeByte(byte(s.bcache))
			s.bcache >>= 8
			s.bcount -= 8
		}
		if unary > 23 {
			s.bcache |= bitMask[23] << s.bcount
			s.bcount += 23
			unary -= 23
		} else {
			s.bcache |= bitMask[unary] << s.bcount
			s.bcount += unary + 1
			unary = 0
		}
		if unary == 0 {
			break
		}
	}
	for s.bcount >= 8 { // put binary
		s.writeByte(byte(s.bcache))
		s.bcache >>= 8
		s.bcount -= 8
	}
	if k != 0 {
		s.bcache |= (outval & bitMask[k]) << s.bcount
		s.bcount += k
	}
}

func (s *fifo) flushBitCache() {
	for s.bcount > 0 {
		s.writeByte(byte(s.bcache))
		s.bcache >>= 8
		if s.bcount > 8 {
			s.bcount -= 8
		} else {
			break
		}
	}
	s.writeCrc32()
}
