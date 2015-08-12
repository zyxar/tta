package tta

import (
	"bytes"
	"testing"
)

func TestUtil(t *testing.T) {
	println("==== util test ====")
}

func TestComputeKeyDigits(t *testing.T) {
	println("[:TestComputeKeyDigits:]")
	var strs = [...]string{"whatisthis?", "1", "", "12", "089q3eoib*(*U(*#$", "~*)(*)@&("}
	var digits = [...][8]byte{
		[8]byte{37, 121, 62, 136, 117, 151, 236, 181},
		[8]byte{90, 13, 77, 214, 205, 142, 75, 114},
		[8]byte{0, 0, 0, 0, 0, 0, 0, 0},
		[8]byte{215, 206, 228, 105, 229, 41, 119, 11},
		[8]byte{35, 135, 48, 205, 86, 61, 214, 216},
		[8]byte{112, 92, 200, 162, 200, 114, 105, 141},
	}
	for i := 0; i < len(strs); i++ {
		b := compute_key_digits([]byte(strs[i]))
		if bytes.Compare(b[:], digits[i][:]) != 0 {
			t.Errorf("compute_key_digits fail @ %d\n", i)
		}
	}
}

func TestConvertPassword(t *testing.T) {
	println("[:TestConvertPassword:]")
	var strs = [...]string{
		"",
		"1",
		"AB",
		"akljsdlfkja;oslduy 98283r7  qiweyr9823475&@^#U#$Y$"}
	var slices = [...][]byte{
		[]byte{},
		[]byte{49},
		[]byte{65, 66},
		[]byte{97, 107, 108, 106, 115, 100, 108, 102, 107, 106, 97, 59, 111, 115, 108, 100, 117, 121, 32, 57, 56, 50, 56, 51, 114, 55, 32, 32, 113, 105, 119, 101, 121, 114, 57, 56, 50, 51, 52, 55, 53, 38, 64, 94, 35, 85, 35, 36, 89, 36}}
	for i := 0; i < len(strs); i++ {
		if bytes.Compare(slices[i], convert_password(strs[i])) != 0 {
			t.Errorf("convert_password fail @ %v\n", strs[i])
		}
	}
}
