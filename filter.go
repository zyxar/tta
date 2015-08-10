// TODO: SSE4 optimization
package tta

func (fs *tta_fltst) init(data [8]byte, shift int32) {
	fs.shift = shift
	fs.round = 1 << uint32(shift-1)
	fs.qm[0] = int32(data[0])
	fs.qm[1] = int32(data[1])
	fs.qm[2] = int32(data[2])
	fs.qm[3] = int32(data[3])
	fs.qm[4] = int32(data[4])
	fs.qm[5] = int32(data[5])
	fs.qm[6] = int32(data[6])
	fs.qm[7] = int32(data[7])
}

func (fs *tta_fltst) hybrid_filter_dec(in *int32) {
	pa := fs.dl
	pb := fs.qm
	pm := fs.dx
	sum := fs.round
	if fs.error < 0 {
		pb[0] -= pm[0]
		pb[1] -= pm[1]
		pb[2] -= pm[2]
		pb[3] -= pm[3]
		pb[4] -= pm[4]
		pb[5] -= pm[5]
		pb[6] -= pm[6]
		pb[7] -= pm[7]
	} else if fs.error > 0 {
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
	fs.error = *in
	*in += (sum >> uint32(fs.shift))
	pa[4] = -pa[5]
	pa[5] = -pa[6]
	pa[6] = *in - pa[7]
	pa[7] = *in
	pa[5] += pa[6]
	pa[4] += pa[5]
}

func (fs *tta_fltst) hybrid_filter_enc(in *int32) {
	pa := fs.dl
	pb := fs.qm
	pm := fs.dx
	sum := fs.round
	if fs.error < 0 {
		pb[0] -= pm[0]
		pb[1] -= pm[1]
		pb[2] -= pm[2]
		pb[3] -= pm[3]
		pb[4] -= pm[4]
		pb[5] -= pm[5]
		pb[6] -= pm[6]
		pb[7] -= pm[7]
	} else if fs.error > 0 {
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

	*in -= (sum >> uint32(fs.shift))
	fs.error = *in
}
