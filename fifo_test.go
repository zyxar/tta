package tta

import (
	"testing"
)

func TestReadByte(t *testing.T) {
	t.Parallel()
	fifo := fifo{}
	for i := 0; i < fifoBufferSize; i++ {
		fifo.buffer[i] = byte(i)
	}
	fifo.pos = 0
	fifo.end = fifoBufferSize
	for i := 0; i < fifoBufferSize; i++ {
		if fifo.readByte() != byte(i&0xFF) {
			t.Errorf("readByte fail @ %d\n", i)
		}
		if fifo.count != uint32(i+1) {
			t.Errorf("readByte fail @ count - %d\n", i)
		}
		if fifo.pos != int32(i+1) {
			t.Errorf("readByte fail @ pos - %d\n", i)
		}
	}
	if fifo.count != fifoBufferSize {
		t.Error("readByte fail @ count")
	}
	if fifo.pos != fifoBufferSize {
		t.Error("readByte fail @ pos")
	}
}

func TestReadUint16(t *testing.T) {
	t.Parallel()
	fifo := fifo{}
	for i := 0; i < fifoBufferSize; i++ {
		fifo.buffer[i] = byte(i)
	}
	fifo.pos = 0
	fifo.end = fifoBufferSize
	var v uint16
	for i := 0; i < fifoBufferSize/2; i++ {
		v = uint16((i*2+1)<<8) | (uint16(i*2) & 0xFF)
		if fifo.readUint16() != v {
			t.Errorf("readUint16 fail @ %d\n", i)
		}
		if fifo.count != uint32(i*2+2) {
			t.Errorf("readUint16 fail @ count - %d\n", i)
		}
		if fifo.pos != int32(i*2+2) {
			t.Errorf("readUint16 fail @ pos - %d\n", i)
		}
	}
	if fifo.count != fifoBufferSize {
		t.Error("readUint16 fail @ count")
	}
	if fifo.pos != fifoBufferSize {
		t.Error("readUint16 fail @ pos")
	}
}

func TestReadUint32(t *testing.T) {
	t.Parallel()
	fifo := fifo{}
	for i := 0; i < fifoBufferSize; i++ {
		fifo.buffer[i] = byte(i)
	}
	fifo.pos = 0
	fifo.end = fifoBufferSize
	var v uint32
	for i := 0; i < fifoBufferSize/4; i++ {
		if fifo.count != uint32(i*4) {
			t.Errorf("readUint32 fail @ count - %d\n", i)
		}
		if fifo.pos != int32(i*4) {
			t.Errorf("readUint32 fail @ pos - %d\n", i)
		}
		v = uint32((i*4+3)<<24&0xFF000000) | uint32((i*4+2)<<16&0xFF0000) | uint32((i*4+1)<<8&0xFF00) | (uint32(i*4) & 0xFF)
		if fifo.readUint32() != v {
			t.Errorf("readUint32 fail @ %d\n", i)
		}
	}
	if fifo.count != fifoBufferSize {
		t.Error("readUint32 fail @ count")
	}
	if fifo.pos != fifoBufferSize {
		t.Error("readUint32 fail @ pos")
	}
}

func TestWriteByte(t *testing.T) {
	t.Parallel()
	fifo := fifo{}
	fifo.pos = 0
	fifo.end = fifoBufferSize
	for i := 0; i < fifoBufferSize; i++ {
		if err := fifo.writeByte(byte(i)); err != nil {
			t.Errorf("writeByte fail @ %d, %v\n", i, err)
		}
		if fifo.count != uint32(i+1) {
			t.Errorf("writeByte fail @ count - %d\n", i)
		}
		if fifo.pos != int32(i+1) {
			t.Errorf("writeByte fail @ pos - %d\n", i)
		}
	}
	if fifo.count != fifoBufferSize {
		t.Error("writeByte fail @ count")
	}
	if fifo.pos != fifoBufferSize {
		t.Error("writeByte fail @ pos")
	}
}
