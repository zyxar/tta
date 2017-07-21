package filter

import (
	"testing"
	"unsafe"
)

func BenchmarkEncodeSSE4(b *testing.B) {
	f := New([8]byte{1, 2, 3, 4, 5, 6, 7, 8}, 8)
	var in int32
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_HybridFilterEncodeSSE4(unsafe.Pointer(f), unsafe.Pointer(&in))
		}
	})
}

func BenchmarkEncodeSSE2(b *testing.B) {
	f := New([8]byte{1, 2, 3, 4, 5, 6, 7, 8}, 8)
	var in int32
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_HybridFilterEncodeSSE2(unsafe.Pointer(f), unsafe.Pointer(&in))
		}
	})
}

func BenchmarkEncodeCompat(b *testing.B) {
	f := New([8]byte{1, 2, 3, 4, 5, 6, 7, 8}, 8)
	var in int32
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_HybridFilterEncodeCompat(unsafe.Pointer(f), unsafe.Pointer(&in))
		}
	})
}

func BenchmarkDecodeSSE4(b *testing.B) {
	f := New([8]byte{1, 2, 3, 4, 5, 6, 7, 8}, 8)
	var in int32
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_HybridFilterDecodeSSE4(unsafe.Pointer(f), unsafe.Pointer(&in))
		}
	})
}

func BenchmarkDecodeSSE2(b *testing.B) {
	f := New([8]byte{1, 2, 3, 4, 5, 6, 7, 8}, 8)
	var in int32
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_HybridFilterDecodeSSE2(unsafe.Pointer(f), unsafe.Pointer(&in))
		}
	})
}

func BenchmarkDecodeCompat(b *testing.B) {
	f := New([8]byte{1, 2, 3, 4, 5, 6, 7, 8}, 8)
	var in int32
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_HybridFilterDecodeCompat(unsafe.Pointer(f), unsafe.Pointer(&in))
		}
	})
}
