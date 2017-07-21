package filter

import (
	"github.com/klauspost/cpuid"
)

const (
	cpuArchUNKNOWN = iota
	cpuArchSSE2
	cpuArchSSE4
)

// CPUArch indicates currently cpu architecture: 0 general; 1 sse2 enabled; 2 sse4 enabled.
var CPUArch = cpuArchUNKNOWN

func init() {
	if cpuid.CPU.SSE4() {
		CPUArch = cpuArchSSE4
	} else if cpuid.CPU.SSE2() {
		CPUArch = cpuArchSSE2
	}
}
