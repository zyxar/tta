package tta

import (
	// "fmt"
	"testing"
)

func TestReadByte(t *testing.T) {
	fifo := tta_fifo{}
	for i := 0; i < TTA_FIFO_BUFFER_SIZE; i++ {
		fifo.buffer[i] = byte(i)
	}
	fifo.pos = 0
	fifo.end = TTA_FIFO_BUFFER_SIZE
	for i := 0; i < TTA_FIFO_BUFFER_SIZE; i++ {
		if fifo.read_byte() != byte(i&0xFF) {
			t.Errorf("read_byte fail @ %d\n", i)
		}
		if fifo.count != uint32(i+1) {
			t.Errorf("read_byte fail @ count - %d\n", i)
		}
		if fifo.pos != int32(i+1) {
			t.Errorf("read_byte fail @ pos - %d\n", i)
		}
	}
	if fifo.count != TTA_FIFO_BUFFER_SIZE {
		t.Error("read_byte fail @ count")
	}
	if fifo.pos != TTA_FIFO_BUFFER_SIZE {
		t.Error("read_byte fail @ pos")
	}
}

func TestReadUint16(t *testing.T) {
	fifo := tta_fifo{}
	for i := 0; i < TTA_FIFO_BUFFER_SIZE; i++ {
		fifo.buffer[i] = byte(i)
	}
	fifo.pos = 0
	fifo.end = TTA_FIFO_BUFFER_SIZE
	var v uint16
	for i := 0; i < TTA_FIFO_BUFFER_SIZE/2; i++ {
		v = uint16((i*2+1)<<8) | (uint16(i*2) & 0xFF)
		if fifo.read_uint16() != v {
			t.Errorf("read_uint16 fail @ %d\n", i)
		}
		if fifo.count != uint32(i*2+2) {
			t.Errorf("read_uint16 fail @ count - %d\n", i)
		}
		if fifo.pos != int32(i*2+2) {
			t.Errorf("read_uint16 fail @ pos - %d\n", i)
		}
	}
	if fifo.count != TTA_FIFO_BUFFER_SIZE {
		t.Error("read_uint16 fail @ count")
	}
	if fifo.pos != TTA_FIFO_BUFFER_SIZE {
		t.Error("read_uint16 fail @ pos")
	}
}

func TestReadUint32(t *testing.T) {
	fifo := tta_fifo{}
	for i := 0; i < TTA_FIFO_BUFFER_SIZE; i++ {
		fifo.buffer[i] = byte(i)
	}
	fifo.pos = 0
	fifo.end = TTA_FIFO_BUFFER_SIZE
	var v uint32
	for i := 0; i < TTA_FIFO_BUFFER_SIZE/4; i++ {
		if fifo.count != uint32(i*4) {
			t.Errorf("read_uint32 fail @ count - %d\n", i)
		}
		if fifo.pos != int32(i*4) {
			t.Errorf("read_uint32 fail @ pos - %d\n", i)
		}
		v = uint32((i*4+3)<<24&0xFF000000) | uint32((i*4+2)<<16&0xFF0000) | uint32((i*4+1)<<8&0xFF00) | (uint32(i*4) & 0xFF)
		if fifo.read_uint32() != v {
			t.Errorf("read_uint32 fail @ %d\n", i)
		}
	}
	if fifo.count != TTA_FIFO_BUFFER_SIZE {
		t.Error("read_uint32 fail @ count")
	}
	if fifo.pos != TTA_FIFO_BUFFER_SIZE {
		t.Error("read_uint32 fail @ pos")
	}
}
