package main

import (
	"log"
	"time"

	"bitbucket.org/gmcbay/i2c"
)

func main() {
	bp, err := i2c.Bus(1)
	if err != nil {
		log.Panicf("failed to create bus: %v\n", err)
	}
	bp.WriteByte(0x70, 0x21, 0x00)
	bp.WriteByte(0x70, 0x81, 0x00)
	bp.WriteByte(0x70, 0xe0, 0x00)

	for t := 0; t < 100; t++ {
		for ix := 0; ix < 16; ix += 2 {
			bp.WriteByte(0x70, byte(ix), 0xcc)
		}
		time.Sleep(200 * time.Millisecond)
		for ix := 0; ix < 16; ix += 2 {
			bp.WriteByte(0x70, byte(ix), 0xff)
		}
		time.Sleep(200 * time.Millisecond)
	}

}
