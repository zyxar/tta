package tta

import (
	"bytes"
	"crypto/sha256"
	"io/ioutil"
	"os"
	"testing"
)

func TestComprehensive(t *testing.T) {
	testfile := "tta_comprehensive_test.wav"
	if info, err := os.Stat(testfile); err == nil {
		if !info.IsDir() {
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
			sum := sha256.Sum256(p)
			tmpfile1 := os.TempDir() + "/tta_comprehensive_test_compressed.tta"
			tmpfile2 := os.TempDir() + "/tta_comprehensive_test_decompressed.wav"
			outfile, err := os.Create(tmpfile1)
			if err != nil {
				t.Fatal(err)
			}
			defer outfile.Close()
			defer os.Remove(tmpfile1)
			if err = Compress(fd, outfile, "", nil); err != nil {
				t.Error(err)
			}
			outfile.Seek(0, os.SEEK_SET)
			outfile2, err := os.Create(tmpfile2)
			if err != nil {
				t.Fatal(err)
			}
			defer outfile2.Close()
			defer os.Remove(tmpfile2)
			if err = Decompress(outfile, outfile2, "", nil); err != nil {
				t.Error(err)
			}
			outfile2.Seek(0, os.SEEK_SET)
			p, err = ioutil.ReadAll(outfile2)
			if err != nil {
				t.Fatal(err)
			}
			sum2 := sha256.Sum256(p)
			if bytes.Compare(sum2[:], sum[:]) != 0 {
				t.Errorf("Checksum fail, expected: %x\n", sum)
			}
		}
	}
}
