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
	defer infile.Close()
	outfile, err := os.Create(os.TempDir() + "/sample_decompressed.wav")
	if err != nil {
		t.Fatal(err)
	}
	defer outfile.Close()
	if err = Decompress(infile, outfile, "", nil); err != nil {
		t.Error(err)
	}
}
