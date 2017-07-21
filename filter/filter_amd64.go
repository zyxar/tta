//+build !noasm !appengine
//+build amd64

package filter

import (
	"unsafe"
)

//go:noescape
func _HybridFilterDecSSE4(in, err, qm, dx, dl unsafe.Pointer, round int32, shift uint32)

//go:noescape
func _HybridFilterEncSSE4(in, err, qm, dx, dl unsafe.Pointer, round int32, shift uint32)

//go:noescape
func _HybridFilterDecSSE2(in, err, qm, dx, dl unsafe.Pointer, round int32, shift uint32)

//go:noescape
func _HybridFilterEncSSE2(in, err, qm, dx, dl unsafe.Pointer, round int32, shift uint32)
