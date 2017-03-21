package main

import (
  "math/rand"
  "time"
	"os"
	"reflect"
	"sync"
	"syscall"
	"unsafe"
	"fmt"
)

const (
	memLength   = 4096
	pinMask uint32 = 7 // 0b111 - pinmode is 3 bits
)

var (
	memlock sync.Mutex
	mem     []uint32
	mem8    []uint8
)

func Open() (err error) {
  file, err := os.OpenFile("/dev/gpiomem", os.O_RDWR|os.O_SYNC, 0)
	if err != nil {
		return
	}
	defer file.Close()

	memlock.Lock()
	defer memlock.Unlock()

	// Memory map GPIO registers to byte array
	mem8, err = syscall.Mmap(
		int(file.Fd()),
		0,
		memLength,
		syscall.PROT_READ|syscall.PROT_WRITE,
		syscall.MAP_SHARED)

	if err != nil {
		return
	}

	// Convert mapped byte memory to unsafe []uint32 pointer, adjust length as needed
	header := *(*reflect.SliceHeader)(unsafe.Pointer(&mem8))
	header.Len /= (32 / 8) // (32 bit = 4 bytes)
	header.Cap /= (32 / 8)

	mem = *(*[]uint32)(unsafe.Pointer(&header))

	return nil
}

// Close unmaps GPIO memory
func Close() error {
	memlock.Lock()
	defer memlock.Unlock()
	return syscall.Munmap(mem8)
}

const (
  p = uint8(4)
  fsel = uint8(p) / 10
  shift = (uint8(p) % 10) * 3
  clearReg = p/32 + 10
  setReg = p/32 + 7
  v uint32 = 1 << (p & 31)
)

const per = 7

func bitLow() {
  for ix := 0; ix < per; ix++ {
    mem[setReg] = v
  }
  for ix := 0; ix < 3*per; ix++ {
    mem[clearReg] = v
  }
}

func bitHigh() {
  for ix := 0; ix < 3*per; ix++ {
    mem[setReg] = v
  }
  for ix := 0; ix < per; ix++ {
    mem[clearReg] = v
  }
}

func bitSet(high bool) {
  if high {
    bitHigh()
  } else {
    bitLow()
  }
}

func byteSet(b uint8) {
 mask := uint8(1) << 7
 for ix := 0; ix < 8; ix++ {
   bitSet((b & mask) > 0)
   mask >>= 1
 }
}

func bytesSet(bs ...uint8) {
  for _, b := range bs {
    byteSet(b)
  }
  mem[clearReg] = v
}

func main() {
	if err := Open(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer Close()

	memlock.Lock()
	defer memlock.Unlock()

  mem[fsel] = (mem[fsel] &^ (pinMask << shift)) | (1 << shift)

	for {
    for ix := 0; ix < 8; ix++ {
      byteSet(uint8(rand.Intn(5)))
      byteSet(uint8(rand.Intn(5)))
      byteSet(uint8(rand.Intn(5)))
    }
//    bytesSet(
//      64, 0, 0,  0, 64, 0,  0, 0, 64,  4, 40, 64,
//      64, 0, 0,  0, 64, 0,  0, 0, 64,  4, 40, 64,
//      64, 0, 0,  0, 64, 0,  0, 0, 64,  4, 40, 64,
//      64, 0, 0,  0, 64, 0,  0, 0, 64,  4, 40, 64)
    mem[clearReg] = v
    time.Sleep(500 * time.Millisecond)
	}
}
