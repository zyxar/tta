package filter

import (
	"testing"
)

func BenchmarkEncodeCompat(b *testing.B) {
	f := New([8]byte{1, 2, 3, 4, 5, 6, 7, 8}, 8)
	var in int32
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			encodeCompat(f, &in)
		}
	})
}

func BenchmarkDecodeCompat(b *testing.B) {
	f := New([8]byte{1, 2, 3, 4, 5, 6, 7, 8}, 8)
	var in int32
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			decodeCompat(f, &in)
		}
	})
}
