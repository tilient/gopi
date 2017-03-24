package main

import (
	"bitbucket.org/gmcbay/i2c"
	"log"
)

func main() {
	bp, err := i2c.Bus(1)
	if err != nil {
		log.Panicf("failed to create bus: %v\n", err)
	}
	bp.WriteByte(0x70, 0x21, 0x00)
	bp.WriteByte(0x70, 0x81, 0x00)
	bp.WriteByte(0x70, 0xe0, 0x00)

	for ix := 0; ix < 16; ix += 2 {
		bp.WriteByte(0x70, byte(ix), 0xaa)
	}
}
