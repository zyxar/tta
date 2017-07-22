package filter

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

type codec func(fs *Filter, in *int32)

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
	decode(f, in)
}

func (f *Filter) Encode(in *int32) {
	encode(f, in)
}

func decodeCompat(f *Filter, in *int32) {
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

func encodeCompat(f *Filter, in *int32) {
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
