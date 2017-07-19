package filter

import (
	"github.com/klauspost/cpuid"
)

var sseEnabled bool

func init() {
	if cpuid.CPU.SSE4() {
		// sseEnabled = true // use this if SSE optimization is done
	}
}
