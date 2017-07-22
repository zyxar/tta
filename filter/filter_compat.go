//+build !amd64

package filter

func DecodeCompat(f *Filter, in *int32) {
	decodeCompat(f, in)
}

func EncodeCompat(f *Filter, in *int32) {
	encodeCompat(f, in)
}
