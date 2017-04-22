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
	bus           *spi.Device
	nrOfMatrices  int
	rowsPerMatrix int
	buffer        []byte
	font          [][]byte
}

func newPixelMatrix(devstr string, brightness byte,
	nrOfMatrices, rowsPerMatrix int,
	font [][]byte) *PixelMatrix {
	bus, err := spi.Open(
		&spi.Devfs{devstr, spi.Mode0, 4000000})
	if err != nil {
		log.Fatal(err)
	}
	buffer := make([]byte, nrOfMatrices*rowsPerMatrix)
	pm := &PixelMatrix{
		bus:           bus,
		nrOfMatrices:  nrOfMatrices,
		rowsPerMatrix: rowsPerMatrix,
		buffer:        buffer,
		font:          font,
	}
	pm.sendCmd(max7219_SCAN_LIMIT, 7)
	pm.sendCmd(max7219_DECODE_MODE, 0)
	pm.sendCmd(max7219_DISPLAY_TEST, 0)
	pm.sendCmd(max7219_SHUTDOWN, 1)
	pm.sendCmd(max7219_INTENSITY, brightness)
	pm.clear()
	pm.flush()
	return pm
}

func (pm *PixelMatrix) Close() {
	pm.bus.Close()
}

// ===========================================================

func (pm *PixelMatrix) flush() {
	buf := make([]byte, 2*pm.nrOfMatrices)
	for line := 0; line < pm.rowsPerMatrix; line++ {
		for matrix := 0; matrix < pm.nrOfMatrices; matrix++ {
			buf[matrix*2] = byte(max7219_DIGIT0 + line)
			buf[matrix*2+1] =
				pm.buffer[matrix*pm.rowsPerMatrix+line]
		}
		pm.bus.Tx(buf[:], nil)
	}
}

// ===========================================================

const (
	max7219_DIGIT0       = 1
	max7219_DECODE_MODE  = 9
	max7219_INTENSITY    = 10
	max7219_SCAN_LIMIT   = 11
	max7219_SHUTDOWN     = 12
	max7219_DISPLAY_TEST = 0x0F
)

func (pm *PixelMatrix) sendCmd(register, value byte) {
	data := []byte{register, value}
	for matrix := 0; matrix < pm.nrOfMatrices; matrix++ {
		pm.bus.Tx(data, nil)
	}
}

// ===========================================================

func (pm *PixelMatrix) clear() {
	for ix := range pm.buffer {
		pm.buffer[ix] = 0
	}
}

func (pm *PixelMatrix) setPixel(x, y int) {
	if x < 0 {
		return
	}
	if y < 0 {
		return
	}
	if y >= pm.rowsPerMatrix {
		return
	}
	if x >= (8 * pm.nrOfMatrices) {
		return
	}
	line := pm.rowsPerMatrix - 1 - y
	matrix := x / 8
	bit := (byte)(1 << (uint)(7-(x%8)))
	pm.buffer[matrix*pm.rowsPerMatrix+line] |= bit
}

// ===========================================================

func (pm *PixelMatrix) plotString(str string, xPos int) int {
	x := xPos
	for _, ch := range str {
		x = pm.plotChar((byte)(ch), x)
		x += 1
	}
	return x
}

func (pm *PixelMatrix) plotChar(ch byte, xPos int) int {
	x := xPos
	bitLines := pm.font[ch-(byte)(' ')]
	for _, bitLine := range bitLines {
		bitPos := (byte)(128)
		for y := 0; y < 8; y++ {
			bit := bitLine & bitPos
			if bit > 0 {
				pm.setPixel(x, y)
			}
			bitPos >>= 1
		}
		x += 1
	}
	return x
}

// ===========================================================
