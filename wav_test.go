package tta

import (
	"bytes"
	"os"
	"testing"
)

var wav_slice = []byte{0x52, 0x49, 0x46, 0x46, 0x98, 0x03, 0x00, 0x00, 0x57, 0x41, 0x56, 0x45, 0x66, 0x6d, 0x74, 0x20,
	0x10, 0x00, 0x00, 0x00, 0x01, 0x00, 0x02, 0x00, 0x44, 0xac, 0x00, 0x00, 0x10, 0xb1, 0x02, 0x00,
	0x04, 0x00, 0x10, 0x00}
var wav_size = uint32(0x0374)

func TestReadHeader(t *testing.T) {
	file, err := os.Open("./data/sample.wav")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	wav := WaveHeader{}
	if size, err := wav.Read(file); err != nil {
		t.Error(err)
	} else {
		if bytes.Compare(wav_slice, wav.Bytes()) != 0 || size != wav_size {
			t.Error("WaveHeader::Read fail")
		}
	}
}

func TestWriteHeader(t *testing.T) {
	filename := os.TempDir() + "/tta_tmp.wav"
	file, err := os.Create(filename)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	defer os.Remove(filename)
	wav := WaveHeader{}
	copy(wav.Bytes(), wav_slice)
	if err = wav.Write(file, wav_size); err != nil {
		t.Error(err)
	}
}
