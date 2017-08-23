//+build amd64

package filter

import (
	"github.com/klauspost/cpuid"
)

func init() {
	if cpuid.CPU.SSE4() {
		encode = SSE4Encode
		decode = SSE4Decode
	} else if cpuid.CPU.SSE2() {
		encode = SSE2Encode
		decode = SSE2Decode
	} else {
		encode = X64Encode
		decode = X64Decode
	}
}
