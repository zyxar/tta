//+build !noasm !appengine
//+build amd64

package filter

import (
	"unsafe"
)

func SSE4Decode(f *Filter, in *int32) { sse4Decode(unsafe.Pointer(f), unsafe.Pointer(in)) }
func SSE4Encode(f *Filter, in *int32) { sse4Encode(unsafe.Pointer(f), unsafe.Pointer(in)) }
func SSE2Decode(f *Filter, in *int32) { sse2Decode(unsafe.Pointer(f), unsafe.Pointer(in)) }
func SSE2Encode(f *Filter, in *int32) { sse2Encode(unsafe.Pointer(f), unsafe.Pointer(in)) }
func X64Decode(f *Filter, in *int32)  { x64Decode(unsafe.Pointer(f), unsafe.Pointer(in)) }
func X64Encode(f *Filter, in *int32)  { x64Encode(unsafe.Pointer(f), unsafe.Pointer(in)) }

//go:noescape
func sse4Decode(fs, in unsafe.Pointer)

//go:noescape
func sse4Encode(fs, in unsafe.Pointer)

//go:noescape
func sse2Decode(fs, in unsafe.Pointer)

//go:noescape
func sse2Encode(fs, in unsafe.Pointer)

//go:noescape
func x64Decode(fs, in unsafe.Pointer)

//go:noescape
func x64Encode(fs, in unsafe.Pointer)
