package filter

import (
	"github.com/klauspost/cpuid"
)

const (
	cpuArchUNKNOWN = iota
	cpuArchSSE2
	cpuArchSSE4
)

var CPUArch = cpuArchUNKNOWN

func init() {
	if cpuid.CPU.SSE4() {
		CPUArch = cpuArchSSE4
	} else if cpuid.CPU.SSE2() {
		CPUArch = cpuArchSSE2
	}
}
