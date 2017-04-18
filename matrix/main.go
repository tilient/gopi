package main

import (
	//"fmt"
	"golang.org/x/exp/io/spi"
	"log"
	"time"
)

func main() {
	device := NewDevice("/dev/spidev0.0", 1)
	defer device.Close()
	font := font()

	for {
		// device.ClearAll()
		// device.plotStringAt(" Wiffel", 0, font())
		// device.Flush()
		// time.Sleep(1 * time.Second)

		// device.ClearAll()
		// device.plotStringAt(" Linda", 0, font())
		// device.Flush()
		// time.Sleep(1 * time.Second)

		device.ClearAll()
		device.plotStringAt(
			time.Now().Format(" 15:04"), 0, font)
		device.Flush()
		time.Sleep(30 * time.Second)
	}
}

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
	if x >= maxX {
		return
	}
	if y >= maxY {
		return
	}
	line := 7 - y
	matId := x / 8
	bit := (byte)(1 << (uint)(7-(x%8)))
	this.buffer[matId*nrOfLines+line] |= bit
}

func (this *Device) plotStringAt(
	str string, xx int, font [][]byte) int {
	x := xx
	for _, ch := range str {
		x = this.plotCharAt((byte)(ch), x, font)
		x += 1
	}
	return x
}

func (this *Device) plotCharAt(ch byte, xx int, font [][]byte) int {
	bits := font[ch-(byte)(' ')]
	x := xx
	for _, bitline := range bits {
		bitpos := (byte)(1) << 7
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

func font() [][]byte {
	return [][]byte{
		{},                                   // ' '
		{0x5F},                               // '!'
		{0x03, 0x00, 0x03},                   // '"'
		{0x24, 0x7E, 0x24, 0x24, 0x7E, 0x24}, // '#'
		{0x2E, 0x2A, 0x7F, 0x2A, 0x3A},       // '$'
		{0x46, 0x26, 0x10, 0x08, 0x64, 0x62}, // '%'
		{0x20, 0x54, 0x4A, 0x54, 0x20, 0x50}, // '&'
		{0x04, 0x02},                         // '''
		{0x3C, 0x42},                         // '('
		{0x42, 0x3C},                         // ')'
		{0x10, 0x54, 0x38, 0x54, 0x10},       // '*'
		{0x10, 0x10, 0x7C, 0x10, 0x10},       // '+'
		{0x80, 0x60},                         // '
		{0x10, 0x10, 0x10, 0x10, 0x10},       // '-'
		{0x60, 0x60},                         // '.'
		{0x40, 0x20, 0x10, 0x08, 0x04},       // '/'
		{0x3C, 0x62, 0x52, 0x4A, 0x46, 0x3C}, // '0'
		{0x44, 0x42, 0x7E, 0x40, 0x40},       // '1'
		{0x64, 0x52, 0x52, 0x52, 0x4C},       // '2'
		{0x24, 0x42, 0x42, 0x4A, 0x4A, 0x34}, // '3'
		{0x30, 0x28, 0x24, 0x7E, 0x20},       // '4'
		{0x2E, 0x4A, 0x4A, 0x4A, 0x32},       // '5'
		{0x3C, 0x4A, 0x4A, 0x4A, 0x30},       // '6'
		{0x02, 0x62, 0x12, 0x0A, 0x06},       // '7'
		{0x34, 0x4A, 0x4A, 0x4A, 0x34},       // '8'
		{0x0C, 0x52, 0x52, 0x52, 0x3C},       // '9'
		{0x48},                                           // ':'
		{0x80, 0x64},                                     // ';'
		{0x10, 0x28, 0x44},                               // '<'
		{0x28, 0x28, 0x28, 0x28, 0x28},                   // '='
		{0x44, 0x28, 0x10},                               // '>'
		{0x04, 0x02, 0x02, 0x52, 0x0A, 0x04},             // '?'
		{0x3C, 0x42, 0x5A, 0x56, 0x5A, 0x1C},             // '@'
		{0x7C, 0x12, 0x12, 0x12, 0x12, 0x7C},             // 'A'
		{0x7E, 0x4A, 0x4A, 0x4A, 0x4A, 0x34},             // 'B'
		{0x3C, 0x42, 0x42, 0x42, 0x42, 0x24},             // 'C'
		{0x7E, 0x42, 0x42, 0x42, 0x24, 0x18},             // 'D'
		{0x7E, 0x4A, 0x4A, 0x4A, 0x4A, 0x42},             // 'E'
		{0x7E, 0x0A, 0x0A, 0x0A, 0x0A, 0x02},             // 'F'
		{0x3C, 0x42, 0x42, 0x52, 0x52, 0x34},             // 'G'
		{0x7E, 0x08, 0x08, 0x08, 0x08, 0x7E},             // 'H'
		{0x42, 0x42, 0x7E, 0x42, 0x42},                   // 'I'
		{0x30, 0x40, 0x40, 0x40, 0x40, 0x3E},             // 'J'
		{0x7E, 0x08, 0x08, 0x14, 0x22, 0x40},             // 'K'
		{0x7E, 0x40, 0x40, 0x40},                         // 'L'
		{0x7E, 0x04, 0x08, 0x08, 0x04, 0x7E},             // 'M'
		{0x7E, 0x04, 0x08, 0x10, 0x20, 0x7E},             // 'N'
		{0x3C, 0x42, 0x42, 0x42, 0x42, 0x3C},             // 'O'
		{0x7E, 0x12, 0x12, 0x12, 0x12, 0x0C},             // 'P'
		{0x3C, 0x42, 0x52, 0x62, 0x42, 0x3C},             // 'Q'
		{0x7E, 0x12, 0x12, 0x12, 0x32, 0x4C},             // 'R'
		{0x24, 0x4A, 0x4A, 0x4A, 0x4A, 0x30},             // 'S'
		{0x02, 0x02, 0x02, 0x7E, 0x02, 0x02, 0x02},       // 'T'
		{0x3E, 0x40, 0x40, 0x40, 0x40, 0x3E},             // 'U'
		{0x1E, 0x20, 0x40, 0x40, 0x20, 0x1E},             // 'V'
		{0x3E, 0x40, 0x20, 0x20, 0x40, 0x3E},             // 'W'
		{0x42, 0x24, 0x18, 0x18, 0x24, 0x42},             // 'X'
		{0x02, 0x04, 0x08, 0x70, 0x08, 0x04, 0x02},       // 'Y'
		{0x42, 0x62, 0x52, 0x4A, 0x46, 0x42},             // 'Z'
		{0x7E, 0x42, 0x42},                               // '['
		{0x04, 0x08, 0x10, 0x20, 0x40},                   // '\'
		{0x42, 0x42, 0x7E},                               // ']'
		{0x08, 0x04, 0x7E, 0x04, 0x08},                   // '^'
		{0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80},       // '_'
		{0x3C, 0x42, 0x99, 0xA5, 0xA5, 0x81, 0x42, 0x3C}, // '`'
		{0x20, 0x54, 0x54, 0x54, 0x78},                   // 'a'
		{0x7E, 0x48, 0x48, 0x48, 0x30},                   // 'b'
		{0x38, 0x44, 0x44, 0x44},                         // 'c'
		{0x30, 0x48, 0x48, 0x48, 0x7E},                   // 'd'
		{0x38, 0x54, 0x54, 0x54, 0x48},                   // 'e'
		{0x7C, 0x0A, 0x02},                               // 'f'
		{0x18, 0xA4, 0xA4, 0xA4, 0xA4, 0x7C},             // 'g'
		{0x7E, 0x08, 0x08, 0x08, 0x70},                   // 'h'
		{0x48, 0x7A, 0x40},                               // 'i'
		{0x40, 0x80, 0x80, 0x7A},                         // 'j'
		{0x7E, 0x18, 0x24, 0x40},                         // 'k'
		{0x3E, 0x40, 0x40},                               // 'l'
		{0x7C, 0x04, 0x78, 0x04, 0x78},                   // 'm'
		{0x7C, 0x04, 0x04, 0x04, 0x78},                   // 'n'
		{0x38, 0x44, 0x44, 0x44, 0x38},                   // 'o'
		{0xFC, 0x24, 0x24, 0x24, 0x18},                   // 'p'
		{0x18, 0x24, 0x24, 0x24, 0xFC, 0x80},             // 'q'
		{0x78, 0x04, 0x04, 0x04},                         // 'r'
		{0x48, 0x54, 0x54, 0x54, 0x20},                   // 's'
		{0x04, 0x3E, 0x44, 0x40},                         // 't'
		{0x3C, 0x40, 0x40, 0x40, 0x3C},                   // 'u'
		{0x0C, 0x30, 0x40, 0x30, 0x0C},                   // 'v'
		{0x3C, 0x40, 0x38, 0x40, 0x3C},                   // 'w'
		{0x44, 0x28, 0x10, 0x28, 0x44},                   // 'x'
		{0x1C, 0xA0, 0xA0, 0xA0, 0x7C},                   // 'y'
		{0x44, 0x64, 0x54, 0x4C, 0x44},                   // 'z'
		{0x08, 0x08, 0x76, 0x42, 0x42},                   // '{'
		{0x7E}, // '|'
		{0x42, 0x42, 0x76, 0x08, 0x08}, // '}'
		{0x04, 0x02, 0x04, 0x02},       // '~'
	}
}
