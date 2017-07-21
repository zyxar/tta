//+build !noasm !appengine
//+build amd64

package filter

import (
	"unsafe"
)

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
