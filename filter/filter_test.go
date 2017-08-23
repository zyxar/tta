package filter

import (
	"testing"
)

func BenchmarkCompatEncode(b *testing.B) {
	f := New([8]byte{1, 2, 3, 4, 5, 6, 7, 8}, 8)
	var in int32
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			CompatEncode(f, &in)
		}
	})
}

func BenchmarkCompatDecode(b *testing.B) {
	f := New([8]byte{1, 2, 3, 4, 5, 6, 7, 8}, 8)
	var in int32
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			CompatDecode(f, &in)
		}
	})
}
