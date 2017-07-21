//+build !noasm !appengine
//+build amd64

package filter

import (
	"unsafe"
)

func DecodeSSE4(f *Filter, in *int32) { _HybridFilterDecodeSSE4(unsafe.Pointer(f), unsafe.Pointer(in)) }
func EncodeSSE4(f *Filter, in *int32) { _HybridFilterEncodeSSE4(unsafe.Pointer(f), unsafe.Pointer(in)) }
func DecodeSSE2(f *Filter, in *int32) { _HybridFilterDecodeSSE2(unsafe.Pointer(f), unsafe.Pointer(in)) }
func EncodeSSE2(f *Filter, in *int32) { _HybridFilterEncodeSSE2(unsafe.Pointer(f), unsafe.Pointer(in)) }
func DecodeCompat(f *Filter, in *int32) {
	_HybridFilterDecodeCompat(unsafe.Pointer(f), unsafe.Pointer(in))
}
func EncodeCompat(f *Filter, in *int32) {
	_HybridFilterEncodeCompat(unsafe.Pointer(f), unsafe.Pointer(in))
}

//go:noescape
func _HybridFilterDecodeSSE4(fs, in unsafe.Pointer)

//go:noescape
func _HybridFilterEncodeSSE4(fs, in unsafe.Pointer)

//go:noescape
func _HybridFilterDecodeSSE2(fs, in unsafe.Pointer)

//go:noescape
func _HybridFilterEncodeSSE2(fs, in unsafe.Pointer)

//go:noescape
func _HybridFilterDecodeCompat(fs, in unsafe.Pointer)

//go:noescape
func _HybridFilterEncodeCompat(fs, in unsafe.Pointer)
