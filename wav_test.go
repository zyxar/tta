package tta

import (
	"bytes"
	"fmt"
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
	wav := WAVE_hdr{}

	if size, err := wav.Read(file); err != nil {
		t.Error(err)
	} else {
		b := wav.toSlice()
		fmt.Printf("sample wav header: %x, %x\n", b, size)
		if bytes.Compare(wav_slice, b) != 0 || size != wav_size {
			t.Error("WAVE_hdr::Read fail")
		}
	}
}

func TestWriteHeader(t *testing.T) {
	wav := WAVE_hdr{}
	filename := os.TempDir() + "/tta_tmp.wav"
	fmt.Println("tmp wave file created at:", filename)
	file, err := os.Create(filename)
	if err != nil {
		t.Fatal(err)
	}
	copy(wav.toSlice(), wav_slice)
	if err = wav.Write(file, wav_size); err != nil {
		t.Error(err)
	}
}