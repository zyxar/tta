package tta

import (
	"bytes"
	"crypto/sha256"
	"io/ioutil"
	"os"
	"testing"
)

func TestComprehensive(t *testing.T) {
	testfile := os.TempDir() + "/tta_comprehensive_test.wav"
	if info, err := os.Stat(testfile); err == nil {
		if !info.IsDir() {
			var sum []byte
			fd, err := os.Open(testfile)
			if err != nil {
				t.Fatal(err)
			}
			defer fd.Close()
			p, err := ioutil.ReadAll(fd)
			if err != nil {
				t.Fatal(err)
			}
			fd.Seek(0, os.SEEK_SET)
			sha := sha256.New()
			sum = sha.Sum(p)
			sha.Reset()
			outfile, err := os.Create(os.TempDir() + "/tta_comprehensive_test_compressed.tta")
			if err != nil {
				t.Fatal(err)
			}
			defer outfile.Close()
			if err = Compress(fd, outfile, "", nil); err != nil {
				t.Error(err)
			}
			outfile.Seek(0, os.SEEK_SET)
			outfile2, err := os.Create(os.TempDir() + "/tta_comprehensive_test_decompressed.wav")
			if err != nil {
				t.Fatal(err)
			}
			defer outfile2.Close()
			if err = Decompress(outfile, outfile2, "", nil); err != nil {
				t.Error(err)
			}
			outfile2.Seek(0, os.SEEK_SET)
			p, err = ioutil.ReadAll(outfile2)
			if err != nil {
				t.Fatal(err)
			}
			if bytes.Compare(sha.Sum(p), sum) != 0 {
				t.Errorf("Checksum fail, expected: %x\n", sum)
			}
		}
	}
}
