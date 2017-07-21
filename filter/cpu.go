package filter

import (
	"github.com/klauspost/cpuid"
)

func init() {
	if cpuid.CPU.SSE4() {
		encode = _HybridFilterEncodeSSE4
		decode = _HybridFilterDecodeSSE4
	} else if cpuid.CPU.SSE2() {
		encode = _HybridFilterEncodeSSE2
		decode = _HybridFilterDecodeSSE2
	} else {
		encode = _HybridFilterEncodeCompat
		decode = _HybridFilterDecodeCompat
	}
}
