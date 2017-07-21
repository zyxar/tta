package filter

import (
	"unsafe"
)

// Filter exposes Decode and Encode methods for data manipulation
type Filter interface {
	Decode(*int32)
	Encode(*int32)
}

type flt struct {
	index  int32
	error  int32
	round  int32
	shift  uint32
	qm     [8]int32
	dx     [24]int32
	dl     [24]int32
	decode func(*int32)
	encode func(*int32)
}

// New creates a Filter based on current CPUArch
func New(data [8]byte, shift uint32) Filter {
	f := flt{}
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
	switch CPUArch {
	case cpuArchSSE4:
		f.decode = f.DecodeSSE4
		f.encode = f.EncodeSSE4
	case cpuArchSSE2:
		f.decode = f.DecodeSSE2
		f.encode = f.EncodeSSE2
	default:
		f.decode = f.DecodeCompat
		f.encode = f.EncodeCompat
	}
	return &f
}

func (f *flt) Decode(in *int32) {
	f.decode(in)
}

func (f *flt) Encode(in *int32) {
	f.encode(in)
}

func (f *flt) DecodeSSE4(in *int32) {
	_HybridFilterDecSSE4(unsafe.Pointer(in), unsafe.Pointer(&f.error), unsafe.Pointer(&f.qm[0]), unsafe.Pointer(&f.dx[0]), unsafe.Pointer(&f.dl[0]), f.round, f.shift)
}

func (f *flt) EncodeSSE4(in *int32) {
	_HybridFilterEncSSE4(unsafe.Pointer(in), unsafe.Pointer(&f.error), unsafe.Pointer(&f.qm[0]), unsafe.Pointer(&f.dx[0]), unsafe.Pointer(&f.dl[0]), f.round, f.shift)
}

func (f *flt) DecodeSSE2(in *int32) {
	_HybridFilterDecSSE2(unsafe.Pointer(in), unsafe.Pointer(&f.error), unsafe.Pointer(&f.qm[0]), unsafe.Pointer(&f.dx[0]), unsafe.Pointer(&f.dl[0]), f.round, f.shift)
}

func (f *flt) EncodeSSE2(in *int32) {
	_HybridFilterEncSSE2(unsafe.Pointer(in), unsafe.Pointer(&f.error), unsafe.Pointer(&f.qm[0]), unsafe.Pointer(&f.dx[0]), unsafe.Pointer(&f.dl[0]), f.round, f.shift)
}

func (f *flt) DecodeCompat(in *int32) {
	pa := f.dl[:]
	pb := f.qm[:]
	pm := f.dx[:]
	sum := f.round
	if f.error < 0 {
		pb[0] -= pm[0]
		pb[1] -= pm[1]
		pb[2] -= pm[2]
		pb[3] -= pm[3]
		pb[4] -= pm[4]
		pb[5] -= pm[5]
		pb[6] -= pm[6]
		pb[7] -= pm[7]
	} else if f.error > 0 {
		pb[0] += pm[0]
		pb[1] += pm[1]
		pb[2] += pm[2]
		pb[3] += pm[3]
		pb[4] += pm[4]
		pb[5] += pm[5]
		pb[6] += pm[6]
		pb[7] += pm[7]
	}
	sum += pa[0]*pb[0] + pa[1]*pb[1] + pa[2]*pb[2] + pa[3]*pb[3] +
		pa[4]*pb[4] + pa[5]*pb[5] + pa[6]*pb[6] + pa[7]*pb[7]

	pm[0] = pm[1]
	pm[1] = pm[2]
	pm[2] = pm[3]
	pm[3] = pm[4]
	pa[0] = pa[1]
	pa[1] = pa[2]
	pa[2] = pa[3]
	pa[3] = pa[4]

	pm[4] = ((pa[4] >> 30) | 1)
	pm[5] = ((pa[5] >> 30) | 2) & ^1
	pm[6] = ((pa[6] >> 30) | 2) & ^1
	pm[7] = ((pa[7] >> 30) | 4) & ^3
	f.error = *in
	*in += (sum >> uint32(f.shift))
	pa[4] = -pa[5]
	pa[5] = -pa[6]
	pa[6] = *in - pa[7]
	pa[7] = *in
	pa[5] += pa[6]
	pa[4] += pa[5]
}

func (f *flt) EncodeCompat(in *int32) {
	pa := f.dl[:]
	pb := f.qm[:]
	pm := f.dx[:]
	sum := f.round
	if f.error < 0 {
		pb[0] -= pm[0]
		pb[1] -= pm[1]
		pb[2] -= pm[2]
		pb[3] -= pm[3]
		pb[4] -= pm[4]
		pb[5] -= pm[5]
		pb[6] -= pm[6]
		pb[7] -= pm[7]
	} else if f.error > 0 {
		pb[0] += pm[0]
		pb[1] += pm[1]
		pb[2] += pm[2]
		pb[3] += pm[3]
		pb[4] += pm[4]
		pb[5] += pm[5]
		pb[6] += pm[6]
		pb[7] += pm[7]
	}

	sum += pa[0]*pb[0] + pa[1]*pb[1] + pa[2]*pb[2] + pa[3]*pb[3] +
		pa[4]*pb[4] + pa[5]*pb[5] + pa[6]*pb[6] + pa[7]*pb[7]

	pm[0] = pm[1]
	pm[1] = pm[2]
	pm[2] = pm[3]
	pm[3] = pm[4]
	pa[0] = pa[1]
	pa[1] = pa[2]
	pa[2] = pa[3]
	pa[3] = pa[4]

	pm[4] = ((pa[4] >> 30) | 1)
	pm[5] = ((pa[5] >> 30) | 2) & ^1
	pm[6] = ((pa[6] >> 30) | 2) & ^1
	pm[7] = ((pa[7] >> 30) | 4) & ^3

	pa[4] = -pa[5]
	pa[5] = -pa[6]
	pa[6] = *in - pa[7]
	pa[7] = *in
	pa[5] += pa[6]
	pa[4] += pa[5]

	*in -= (sum >> uint32(f.shift))
	f.error = *in
}
