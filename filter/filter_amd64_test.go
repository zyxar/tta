//+build !noasm !appengine
//+build amd64

package filter

import (
	"testing"
)

func BenchmarkSSE4Encode(b *testing.B) {
	f := New([8]byte{1, 2, 3, 4, 5, 6, 7, 8}, 8)
	var in int32
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			SSE4Encode(f, &in)
		}
	})
}

func BenchmarkSSE2Encode(b *testing.B) {
	f := New([8]byte{1, 2, 3, 4, 5, 6, 7, 8}, 8)
	var in int32
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			SSE2Encode(f, &in)
		}
	})
}

func BenchmarkX64Encode(b *testing.B) {
	f := New([8]byte{1, 2, 3, 4, 5, 6, 7, 8}, 8)
	var in int32
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			X64Encode(f, &in)
		}
	})
}

func BenchmarkSSE4Decode(b *testing.B) {
	f := New([8]byte{1, 2, 3, 4, 5, 6, 7, 8}, 8)
	var in int32
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			SSE4Decode(f, &in)
		}
	})
}

func BenchmarkSSE2Decode(b *testing.B) {
	f := New([8]byte{1, 2, 3, 4, 5, 6, 7, 8}, 8)
	var in int32
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			SSE2Decode(f, &in)
		}
	})
}

func BenchmarkX64Decode(b *testing.B) {
	f := New([8]byte{1, 2, 3, 4, 5, 6, 7, 8}, 8)
	var in int32
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			X64Decode(f, &in)
		}
	})
}
