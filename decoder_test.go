package tta

import (
	"fmt"
	"os"
	"testing"
)

func TestDecoder(t *testing.T) {
	println("==== decoder test ====")
}

func TestDecompress(t *testing.T) {
	infile, err := os.Open("./data/sample.tta")
	if err != nil {
		t.Fatal(err)
	}
	outfile, err := os.Create(os.TempDir() + "/sample_decompressed.wav")
	if err != nil {
		t.Fatal(err)
	}
	if err = Decompress(infile, outfile, "", func(rate, fnum, frames uint32) {
		pcnt := uint32(float32(fnum) * 100 / float32(frames))
		if (pcnt % 10) == 0 {
			fmt.Printf("\rProgress: %02d%% [%02d]", pcnt, rate)
		}
	}); err != nil {
		t.Error(err)
	}
	println()
}
