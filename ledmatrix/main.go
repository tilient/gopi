package main

import (
	"log"
	"time"

	"bitbucket.org/gmcbay/i2c"
	//"bitbucket.org/gmcbay/i2c/HT16K33"
)

func main() {
	bp, err := i2c.Bus(1)
	if err != nil {
		log.Panicf("failed to create bus: %v\n", err)
	}
	bp.WriteByte(0x70, 0x21, 0x00)
	bp.WriteByte(0x70, 0x81, 0x00)
	bp.WriteByte(0x70, 0xe0, 0x00)

	//devices, err := HT16K33.ParseDevices("0x71:1,0x70:1")
	//if err != nil {
	//		log.Panicf("failed to parse devices: %v\n", err)
	//}
	for {
		bp.WriteByte(0x70, 0xef, 0x00)
		bp.WriteByte(0x71, 0xef, 0x00)

		//HT16K33.ScrollMessage("wiffel & linda",
		//  devices, 200)

		for t := 0; t < 10; t++ {
			for ix := 0; ix < 16; ix += 2 {
				bp.WriteByte(0x70, byte(ix), 0xcc)
				bp.WriteByte(0x71, byte(ix), 0xcc)
			}
			time.Sleep(80 * time.Millisecond)
			for ix := 0; ix < 16; ix += 2 {
				bp.WriteByte(0x70, byte(ix), 0xff)
				bp.WriteByte(0x71, byte(ix), 0xff)
			}
			time.Sleep(80 * time.Millisecond)
		}
		for t := 0; t < 10; t++ {
			bp.WriteByte(0x70, 0xe0, 0x00)
			bp.WriteByte(0x71, 0xe0, 0x00)
			time.Sleep(80 * time.Millisecond)
			bp.WriteByte(0x70, 0xef, 0x00)
			bp.WriteByte(0x71, 0xef, 0x00)
			time.Sleep(80 * time.Millisecond)
		}
		bp.WriteByte(0x70, 0xe0, 0x00)
		bp.WriteByte(0x71, 0xe0, 0x00)
	}

}
