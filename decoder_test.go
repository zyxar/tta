package tta

import (
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
	if err = Decompress(infile, outfile, "", nil); err != nil {
		t.Error(err)
	}
}
