package tta

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
