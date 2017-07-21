package filter

import (
	"testing"
)

func newFlt() *flt {
	t := flt{}
	t.shift = 8
	t.round = 1 << uint32(t.shift-1)
	t.qm = [8]int32{1, 2, 3, 4, 5, 6, 7, 8}
	return &t
}

func BenchmarkEncodeSSE4(b *testing.B) {
	f := newFlt()
	var in int32
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			f.EncodeSSE4(&in)
		}
	})
}

func BenchmarkEncodeSSE2(b *testing.B) {
	f := newFlt()
	var in int32
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			f.EncodeSSE2(&in)
		}
	})
}

func BenchmarkEncodeCompat(b *testing.B) {
	f := newFlt()
	var in int32
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			f.EncodeCompat(&in)
		}
	})
}

func BenchmarkDecodeSSE4(b *testing.B) {
	f := newFlt()
	var in int32
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			f.DecodeSSE4(&in)
		}
	})
}

func BenchmarkDecodeSSE2(b *testing.B) {
	f := newFlt()
	var in int32
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			f.DecodeSSE2(&in)
		}
	})
}

func BenchmarkDecodeCompat(b *testing.B) {
	f := newFlt()
	var in int32
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			f.DecodeCompat(&in)
		}
	})
}
