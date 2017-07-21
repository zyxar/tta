package filter

import (
	"github.com/klauspost/cpuid"
)

func init() {
	if cpuid.CPU.SSE4() {
		encode = EncodeSSE4
		decode = DecodeSSE4
	} else if cpuid.CPU.SSE2() {
		encode = EncodeSSE2
		decode = DecodeSSE2
	} else {
		encode = EncodeCompat
		decode = DecodeCompat
	}
}
