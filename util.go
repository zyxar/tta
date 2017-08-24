package tta

import (
	"encoding/binary"
)

var le = binary.LittleEndian

func computeKeyDigits(p []byte) [8]byte {
	var crcLow, crcHigh uint32 = 0xFFFFFFFF, 0xFFFFFFFF
	for i := 0; i < len(p); i++ {
		index := (crcHigh >> 24) ^ uint32(p[i])&0xFF
		crcHigh = crc64TableHigh[index] ^ ((crcHigh << 8) | (crcLow >> 24))
		crcLow = crc64TableLow[index] ^ (crcLow << 8)
	}
	crcLow ^= 0xFFFFFFFF
	crcHigh ^= 0xFFFFFFFF
	return [8]byte{
		byte((crcLow) & 0xFF),
		byte((crcLow >> 8) & 0xFF),
		byte((crcLow >> 16) & 0xFF),
		byte((crcLow >> 24) & 0xFF),
		byte((crcHigh) & 0xFF),
		byte((crcHigh >> 8) & 0xFF),
		byte((crcHigh >> 16) & 0xFF),
		byte((crcHigh >> 24) & 0xFF),
	}
}

func convertPassword(src string) []byte {
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

func writeBuffer(src int32, p []byte, depth uint32) {
	switch depth {
	case 2:
		le.PutUint16(p, uint16(0xFFFF&src))
	case 1:
		p[0] = byte(0xFF & src)
	case 3:
		le.PutUint16(p, uint16(0xFFFF&src))
		p[2] = byte(0xFF & (src >> 16))
	default:
		le.PutUint32(p, uint32(src))
	}
}

func readBuffer(p []byte, depth uint32) (v int32) {
	switch depth {
	case 2:
		v = int32(int16(le.Uint16(p)))
	case 1:
		v = int32(int8(p[0]))
	case 3:
		v = int32(le.Uint16(p)) | int32(uint32(p[2])<<16)
	default:
		v = int32(le.Uint32(p))
	}
	return
}
