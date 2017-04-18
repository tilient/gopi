package main

import (
	//"fmt"
	"golang.org/x/exp/io/spi"
	"log"
	//"time"
)

const (
	MAX7219_REG_NOOP   byte = 0
	MAX7219_REG_DIGIT0      = iota
	MAX7219_REG_DIGIT1
	MAX7219_REG_DIGIT2
	MAX7219_REG_DIGIT3
	MAX7219_REG_DIGIT4
	MAX7219_REG_DIGIT5
	MAX7219_REG_DIGIT6
	MAX7219_REG_DIGIT7
	MAX7219_REG_DECODEMODE
	MAX7219_REG_INTENSITY
	MAX7219_REG_SCANLIMIT
	MAX7219_REG_SHUTDOWN
	MAX7219_REG_DISPLAYTEST = 0x0F
)

const (
	nrOfLines    = 8
	nrOfMatrices = 4

	maxX = 8 * 4
	maxY = 8
)

type Device struct {
	buffer [nrOfLines * nrOfMatrices]byte
	spi    *spi.Device
}

func NewDevice(devstr string, brightness byte) *Device {
	spi, err := spi.Open(&spi.Devfs{devstr, spi.Mode0, 4000000})
	if err != nil {
		log.Fatal(err)
	}
	this := &Device{spi: spi}
	this.SendCmd(MAX7219_REG_SCANLIMIT, 7)
	this.SendCmd(MAX7219_REG_DECODEMODE, 0)
	this.SendCmd(MAX7219_REG_DISPLAYTEST, 0)
	this.SendCmd(MAX7219_REG_SHUTDOWN, 1)
	this.SendCmd(MAX7219_REG_INTENSITY, brightness)
	this.Flush()
	return this
}

func (this *Device) Close() {
	this.spi.Close()
}

func (this *Device) SendCmd(c, b byte) {
	for matId := 0; matId < nrOfMatrices; matId++ {
		this.spi.Tx([]byte{c, b}, nil)
	}
}

func (this *Device) SetBufferLine(
	matId int, line int, value byte) {
	this.buffer[matId*nrOfLines+line] = value
}

func (this *Device) Flush() {
	for line := 0; line < nrOfLines; line++ {
		var buf [nrOfMatrices * 2]byte
		for matId := 0; matId < nrOfMatrices; matId++ {
			buf[matId*2] = byte(MAX7219_REG_DIGIT0 + line)
			buf[matId*2+1] = this.buffer[matId*nrOfLines+line]
		}
		this.spi.Tx(buf[:], nil)
	}
}

func (this *Device) ClearAll() {
	for ix := range this.buffer {
		this.buffer[ix] = 0
	}
}

func (this *Device) setPixel(x, y int) {
	line := 7 - y
	matId := x / 8
	bit := (byte)(1 << (uint)(7-(x%8)))
	this.buffer[matId*nrOfLines+line] |= bit
}

func main() {
	device := NewDevice("/dev/spidev0.0", 1)
	defer device.Close()

	for x := 0; x < maxX; x++ {
		device.setPixel(x, 3)
		device.setPixel(x, 4)
	}
	for y := 0; y < maxY; y++ {
		device.setPixel(3, y)
		device.setPixel(4, y)
	}
	device.Flush()
	// 	time.Sleep(5000 * time.Millisecond)
	//
	// 	var b byte = 0x55
	// 	for t := 1; t < 20; t++ {
	// 		for line := 0; line < 1; line++ {
	// 			for matId := 0; matId < 1; matId++ {
	// 				device.SetBufferLine(matId, line, b)
	// 			}
	// 			b ^= 0xFF
	// 		}
	// 		device.Flush()
	// 		b ^= 0xFF
	// 		time.Sleep(500 * time.Millisecond)
	// 	}
}
