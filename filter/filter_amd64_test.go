//+build !noasm !appengine
//+build amd64

package filter

import (
	"testing"
)

func BenchmarkEncodeSSE4(b *testing.B) {
	f := New([8]byte{1, 2, 3, 4, 5, 6, 7, 8}, 8)
	var in int32
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			EncodeSSE4(f, &in)
		}
	})
}

func BenchmarkEncodeSSE2(b *testing.B) {
	f := New([8]byte{1, 2, 3, 4, 5, 6, 7, 8}, 8)
	var in int32
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			EncodeSSE2(f, &in)
		}
	})
}

func BenchmarkEncodeCompatX64(b *testing.B) {
	f := New([8]byte{1, 2, 3, 4, 5, 6, 7, 8}, 8)
	var in int32
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			EncodeCompat(f, &in)
		}
	})
}

func BenchmarkDecodeSSE4(b *testing.B) {
	f := New([8]byte{1, 2, 3, 4, 5, 6, 7, 8}, 8)
	var in int32
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			DecodeSSE4(f, &in)
		}
	})
}

func BenchmarkDecodeSSE2(b *testing.B) {
	f := New([8]byte{1, 2, 3, 4, 5, 6, 7, 8}, 8)
	var in int32
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			DecodeSSE2(f, &in)
		}
	})
}

func BenchmarkDecodeCompatX64(b *testing.B) {
	f := New([8]byte{1, 2, 3, 4, 5, 6, 7, 8}, 8)
	var in int32
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			DecodeCompat(f, &in)
		}
	})
}
