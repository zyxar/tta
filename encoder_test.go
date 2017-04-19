package tta

import (
	"os"
	"testing"
)

func TestCompress(t *testing.T) {
	t.Parallel()
	infile, err := os.Open("./data/sample.wav")
	if err != nil {
		t.Fatal(err)
	}
	defer infile.Close()
	outfile, err := os.Create(os.TempDir() + "/sample_compressed.tta")
	if err != nil {
		t.Fatal(err)
	}
	defer outfile.Close()
	defer os.Remove(os.TempDir() + "/sample_compressed.tta")
	if err = Compress(infile, outfile, "", nil); err != nil {
		t.Error(err)
	}
}
