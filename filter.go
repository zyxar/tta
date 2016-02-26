package tta

// TODO: SSE4 optimization

func NewCompatibleFilter(data [8]byte, shift int32) Filter {
	t := ttaFilterCompat{}
	t.shift = shift
	t.round = 1 << uint32(shift-1)
	t.qm[0] = int32(int8(data[0]))
	t.qm[1] = int32(int8(data[1]))
	t.qm[2] = int32(int8(data[2]))
	t.qm[3] = int32(int8(data[3]))
	t.qm[4] = int32(int8(data[4]))
	t.qm[5] = int32(int8(data[5]))
	t.qm[6] = int32(int8(data[6]))
	t.qm[7] = int32(int8(data[7]))
	return &t
}

func (t *ttaFilterCompat) Decode(in *int32) {
	pa := t.dl[:]
	pb := t.qm[:]
	pm := t.dx[:]
	sum := t.round
	if t.error < 0 {
		pb[0] -= pm[0]
		pb[1] -= pm[1]
		pb[2] -= pm[2]
		pb[3] -= pm[3]
		pb[4] -= pm[4]
		pb[5] -= pm[5]
		pb[6] -= pm[6]
		pb[7] -= pm[7]
	} else if t.error > 0 {
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
	t.error = *in
	*in += (sum >> uint32(t.shift))
	pa[4] = -pa[5]
	pa[5] = -pa[6]
	pa[6] = *in - pa[7]
	pa[7] = *in
	pa[5] += pa[6]
	pa[4] += pa[5]
}

func (t *ttaFilterCompat) Encode(in *int32) {
	pa := t.dl[:]
	pb := t.qm[:]
	pm := t.dx[:]
	sum := t.round
	if t.error < 0 {
		pb[0] -= pm[0]
		pb[1] -= pm[1]
		pb[2] -= pm[2]
		pb[3] -= pm[3]
		pb[4] -= pm[4]
		pb[5] -= pm[5]
		pb[6] -= pm[6]
		pb[7] -= pm[7]
	} else if t.error > 0 {
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

	*in -= (sum >> uint32(t.shift))
	t.error = *in
}

func NewSSEFilter(data [8]byte, shift int32) Filter {
	t := ttaFilterSse{}
	t.shift = shift
	t.round = 1 << uint32(shift-1)
	t.qm[0] = int32(int8(data[0]))
	t.qm[1] = int32(int8(data[1]))
	t.qm[2] = int32(int8(data[2]))
	t.qm[3] = int32(int8(data[3]))
	t.qm[4] = int32(int8(data[4]))
	t.qm[5] = int32(int8(data[5]))
	t.qm[6] = int32(int8(data[6]))
	t.qm[7] = int32(int8(data[7]))
	return &t
}

func (t *ttaFilterSse) Decode(in *int32) {
}

func (t *ttaFilterSse) Encode(in *int32) {
}
