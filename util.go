package tta

import (
	"encoding/binary"
)

func compute_key_digits(p []byte) [8]byte {
	var crc_lo, crc_hi uint32 = 0xFFFFFFFF, 0xFFFFFFFF
	for i := 0; i < len(p); i++ {
		index := (crc_hi >> 24) ^ uint32(p[i])&0xFF
		crc_hi = crc64_table_hi[index] ^ ((crc_hi << 8) | (crc_lo >> 24))
		crc_lo = crc64_table_lo[index] ^ (crc_lo << 8)
	}
	crc_lo ^= 0xFFFFFFFF
	crc_hi ^= 0xFFFFFFFF
	return [8]byte{
		byte((crc_lo) & 0xFF),
		byte((crc_lo >> 8) & 0xFF),
		byte((crc_lo >> 16) & 0xFF),
		byte((crc_lo >> 24) & 0xFF),
		byte((crc_hi) & 0xFF),
		byte((crc_hi >> 8) & 0xFF),
		byte((crc_hi >> 16) & 0xFF),
		byte((crc_hi >> 24) & 0xFF),
	}
}

func convert_password(src string) []byte {
	dst := make([]byte, len(src))
	for i := 0; i < len(src); i++ {
		if src[i]&0xF0 == 0xF0 {
			dst[i] = src[i] & 0x0F
		} else if src[i]&0xE0 == 0xE0 {
			dst[i] = src[i] & 0x1F
		} else if src[i]&0xC0 == 0xC0 {
			dst[i] = src[i] & 0x3F
		} else if src[i]&0x80 == 0x80 {
			dst[i] = src[i] & 0x7F
		} else {
			dst[i] = src[i]
		}
	}
	return dst
}

func write_buffer(src int32, p []byte, depth uint32) {
	switch depth {
	case 2:
		binary.LittleEndian.PutUint16(p, uint16(0xFFFF&src))
	case 1:
		p[0] = byte(0xFF & src)
	default:
		binary.LittleEndian.PutUint32(p, uint32(0xFFFF&src))
	}
}

func read_buffer(p []byte, depth uint32) (v int32) {
	switch depth {
	case 2:
		v = int32(int16(binary.LittleEndian.Uint16(p)))
	case 1:
		v = int32(int8(p[0]))
	default:
		v = int32(binary.LittleEndian.Uint32(p))
	}
	return
}
