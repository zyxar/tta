//+build !noasm !appengine
//+build amd64

package filter

import (
	"unsafe"
)

//go:noescape
func __hybrid_filter_dec_sse4(in, err, qm, dx, dl unsafe.Pointer, round int32, shift uint32)

//go:noescape
func __hybrid_filter_enc_sse4(in, err, qm, dx, dl unsafe.Pointer, round int32, shift uint32)

//go:noescape
func __hybrid_filter_dec_sse2(in, err, qm, dx, dl unsafe.Pointer, round int32, shift uint32)

//go:noescape
func __hybrid_filter_enc_sse2(in, err, qm, dx, dl unsafe.Pointer, round int32, shift uint32)
