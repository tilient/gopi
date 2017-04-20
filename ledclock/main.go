package main

import (
	"log"
	"time"

	"golang.org/x/exp/io/spi"
)

// ===========================================================
// Clock
// ===========================================================

func main() {
	pixelMatrix := newPixelMatrix(
		"/dev/spidev0.0", 1, 4, 8, cp437Font())
	defer pixelMatrix.Close()

	for {
		pixelMatrix.clear()
		pixelMatrix.plotString(time.Now().Format("15:04"), 1)
		pixelMatrix.flush()
		time.Sleep(30 * time.Second)
	}
}

// ===========================================================
// Cascaded Max7219 PixelMatrix
// ===========================================================

type PixelMatrix struct {
	spi               *spi.Device
	nrOfMatrices      int
	nrOfRowsPerMatrix int
	buffer            []byte
	font              [][]byte
}

func newPixelMatrix(devstr string, brightness byte,
	nrOfMatrices, nrOfRowsPerMatrix int,
	font [][]byte) *PixelMatrix {
	spi, err := spi.Open(
		&spi.Devfs{devstr, spi.Mode0, 4000000})
	if err != nil {
		log.Fatal(err)
	}
	buffer := make([]byte, nrOfMatrices*nrOfRowsPerMatrix)
	this := &PixelMatrix{
		spi:               spi,
		nrOfMatrices:      nrOfMatrices,
		nrOfRowsPerMatrix: nrOfRowsPerMatrix,
		buffer:            buffer,
		font:              font,
	}
	this.sendCmd(MAX7219_REG_SCANLIMIT, 7)
	this.sendCmd(MAX7219_REG_DECODEMODE, 0)
	this.sendCmd(MAX7219_REG_DISPLAYTEST, 0)
	this.sendCmd(MAX7219_REG_SHUTDOWN, 1)
	this.sendCmd(MAX7219_REG_INTENSITY, brightness)
	this.clear()
	this.flush()
	return this
}

func (this *PixelMatrix) Close() {
	this.spi.Close()
}

// ===========================================================

func (this *PixelMatrix) flush() {
	buf := make([]byte, 2*this.nrOfMatrices)
	for line := 0; line < this.nrOfRowsPerMatrix; line++ {
		for matrix := 0; matrix < this.nrOfMatrices; matrix++ {
			buf[matrix*2] = byte(MAX7219_REG_DIGIT0 + line)
			buf[matrix*2+1] =
				this.buffer[matrix*this.nrOfRowsPerMatrix+line]
		}
		this.spi.Tx(buf[:], nil)
	}
}

// ===========================================================

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

func (this *PixelMatrix) sendCmd(register, value byte) {
	data := []byte{register, value}
	for matrix := 0; matrix < this.nrOfMatrices; matrix++ {
		this.spi.Tx(data, nil)
	}
}

// ===========================================================

func (this *PixelMatrix) clear() {
	for ix := range this.buffer {
		this.buffer[ix] = 0
	}
}

func (this *PixelMatrix) setPixel(x, y int) {
	if x < 0 {
		return
	}
	if y < 0 {
		return
	}
	if y >= this.nrOfRowsPerMatrix {
		return
	}
	if x >= (8 * this.nrOfMatrices) {
		return
	}
	line := this.nrOfRowsPerMatrix - 1 - y
	matrix := x / 8
	bit := (byte)(1 << (uint)(7-(x%8)))
	this.buffer[matrix*this.nrOfRowsPerMatrix+line] |= bit
}

// ===========================================================

func (this *PixelMatrix) plotString(str string, xPos int) int {
	x := xPos
	for _, ch := range str {
		x = this.plotChar((byte)(ch), x)
		x += 1
	}
	return x
}

func (this *PixelMatrix) plotChar(ch byte, xPos int) int {
	x := xPos
	bitlines := this.font[ch-(byte)(' ')]
	for _, bitline := range bitlines {
		bitpos := (byte)(128)
		for y := 0; y < 8; y++ {
			bit := bitline & bitpos
			if bit > 0 {
				this.setPixel(x, y)
			}
			bitpos >>= 1
		}
		x += 1
	}
	return x
}

// ===========================================================
