package tta

import (
	"bytes"
	"testing"
)

func TestComputeKeyDigits(t *testing.T) {
	var strs = [...]string{"whatisthis?", "1", "", "12", "089q3eoib*(*U(*#$", "~*)(*)@&("}
	var digits = [...][8]byte{
		{37, 121, 62, 136, 117, 151, 236, 181},
		{90, 13, 77, 214, 205, 142, 75, 114},
		{0, 0, 0, 0, 0, 0, 0, 0},
		{215, 206, 228, 105, 229, 41, 119, 11},
		{35, 135, 48, 205, 86, 61, 214, 216},
		{112, 92, 200, 162, 200, 114, 105, 141},
	}
	for i := 0; i < len(strs); i++ {
		b := computeKeyDigits([]byte(strs[i]))
		if bytes.Compare(b[:], digits[i][:]) != 0 {
			t.Errorf("computeKeyDigits fail @ %d\n", i)
		}
	}
}

func TestConvertPassword(t *testing.T) {
	var strs = [...]string{
		"",
		"1",
		"AB",
		"akljsdlfkja;oslduy 98283r7  qiweyr9823475&@^#U#$Y$"}
	var slices = [...][]byte{
		{},
		{49},
		{65, 66},
		{97, 107, 108, 106, 115, 100, 108, 102, 107, 106, 97, 59, 111, 115, 108, 100, 117, 121, 32, 57, 56, 50, 56, 51, 114, 55, 32, 32, 113, 105, 119, 101, 121, 114, 57, 56, 50, 51, 52, 55, 53, 38, 64, 94, 35, 85, 35, 36, 89, 36}}
	for i := 0; i < len(strs); i++ {
		if bytes.Compare(slices[i], convertPassword(strs[i])) != 0 {
			t.Errorf("convertPassword fail @ %v\n", strs[i])
		}
	}
}
