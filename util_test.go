package tta

import (
	"bytes"
	"testing"
)

func TestComputeKeyDigits(t *testing.T) {
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
