package filter

import (
	"unsafe"
)

// Filter exposes Decode and Encode methods for data manipulation
type Filter struct {
	index int32
	error int32
	round int32
	shift uint32
	qm    [8]int32
	dx    [24]int32
	dl    [24]int32
}

type codec func(fs, in unsafe.Pointer)

var (
	decode codec
	encode codec
)

// New creates a Filter based on data and shift
func New(data [8]byte, shift uint32) *Filter {
	f := Filter{}
	f.shift = shift
	f.round = 1 << uint32(shift-1)
	f.qm[0] = int32(int8(data[0]))
	f.qm[1] = int32(int8(data[1]))
	f.qm[2] = int32(int8(data[2]))
	f.qm[3] = int32(int8(data[3]))
	f.qm[4] = int32(int8(data[4]))
	f.qm[5] = int32(int8(data[5]))
	f.qm[6] = int32(int8(data[6]))
	f.qm[7] = int32(int8(data[7]))
	return &f
}

func (f *Filter) Decode(in *int32) {
	decode(unsafe.Pointer(f), unsafe.Pointer(in))
}

func (f *Filter) Encode(in *int32) {
	encode(unsafe.Pointer(f), unsafe.Pointer(in))
}
