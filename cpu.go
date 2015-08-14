package tta

import (
	"github.com/klauspost/cpuid"
)

func init() {
	if cpuid.CPU.SSE4() {
		// SSE_Enabled = true // use this if SSE optimization is done
	}
}
