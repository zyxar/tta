package wave

import (
	"bytes"
	"os"
	"testing"
)

var wavSlice = []byte{0x52, 0x49, 0x46, 0x46, 0x98, 0x03, 0x00, 0x00, 0x57, 0x41, 0x56, 0x45, 0x66, 0x6d, 0x74, 0x20,
	0x10, 0x00, 0x00, 0x00, 0x01, 0x00, 0x02, 0x00, 0x44, 0xac, 0x00, 0x00, 0x10, 0xb1, 0x02, 0x00,
	0x04, 0x00, 0x10, 0x00}
var wavSize = uint32(0x0374)

func TestReadHeader(t *testing.T) {
	t.Parallel()
	file, err := os.Open("../data/sample.wav")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	wav := Header{}
	if size, err := wav.Read(file); err != nil {
		t.Error(err)
	} else {
		if bytes.Compare(wavSlice, wav.Bytes()) != 0 || size != wavSize {
			t.Error("Header::Read fail")
		}
	}
}

func TestWriteHeader(t *testing.T) {
	t.Parallel()
	filename := os.TempDir() + "/tta_tmp.wav"
	file, err := os.Create(filename)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	defer os.Remove(filename)
	wav := Header{}
	copy(wav.Bytes(), wavSlice)
	if err = wav.Write(file, wavSize); err != nil {
		t.Error(err)
	}
}
