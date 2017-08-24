package wave

import (
	"bytes"
	"os"
	"testing"
)

var wavSlice = []byte{
	0x52, 0x49, 0x46, 0x46, 0x98, 0x03, 0x00, 0x00, 0x57, 0x41, 0x56, 0x45, 0x66, 0x6d, 0x74, 0x20,
	0x10, 0x00, 0x00, 0x00, 0x01, 0x00, 0x02, 0x00, 0x44, 0xac, 0x00, 0x00, 0x10, 0xb1, 0x02, 0x00,
	0x04, 0x00, 0x10, 0x00}

func TestReadHeader(t *testing.T) {
	t.Parallel()
	file, err := os.Open("../data/sample.wav")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	if wav, size, err := ReadHeader(file); err != nil {
		t.Error(err)
	} else {
		if bytes.Compare(wavSlice, wav.Bytes()) != 0 || size != wav.ChunkSize-szHeader {
			t.Error("Header::Read fail")
		}
	}
}

func TestWriteHeader(t *testing.T) {
	t.Parallel()
	var buf bytes.Buffer
	wav := Header{}
	copy(wav.Bytes(), wavSlice)
	if n, err := wav.WriteTo(&buf); err != nil {
		t.Error(err)
	} else if n != int64(szHeader+szSubchunkHeader) {
		t.Errorf("write %d bytes, expected %d bytes", n, szHeader+szSubchunkHeader)
	}
}
