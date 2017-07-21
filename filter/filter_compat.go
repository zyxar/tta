//+build !amd64

package filter

func DecodeCompat(f *Filter, in *int32) {
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

func EncodeCompat(f *Filter, in *int32) {
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
