package tta

import (
	"io"
)

type Info struct {
	format  uint32 // audio format
	nch     uint32 // number of channels
	bps     uint32 // bits per sample
	sps     uint32 // samplerate (sps)
	samples uint32 // data length in samples
}

type Filter interface {
	Decode(*int32)
	Encode(*int32)
}

type ttaFilterCompat struct {
	index int32
	error int32
	round int32
	shift int32
	qm    [8]int32
	dx    [24]int32
	dl    [24]int32
}

type ttaFilterSse ttaFilterCompat

type ttaAdapt struct {
	k0   uint32
	k1   uint32
	sum0 uint32
	sum1 uint32
}

func (rice *ttaAdapt) init(k0, k1 uint32) {
	rice.k0 = k0
	rice.k1 = k1
	rice.sum0 = shift16[k0]
	rice.sum1 = shift16[k1]
}

type ttaCodec struct {
	filter Filter
	rice   ttaAdapt
	prev   int32
}

type ttaFifo struct {
	buffer [fifoBufferSize]byte
	pos    int32
	end    int32
	bcount uint32 // count of bits in cache
	bcache uint32 // bit cache
	crc    uint32
	count  uint32
	io     io.ReadWriteSeeker
}

type Callback func(uint32, uint32, uint32)
