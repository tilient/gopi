package main

import (
	"fmt"
	"os"

	"github.com/mrmorphic/hwio"
)

func main() {
	ledPin, err := hwio.GetPin("gpio4")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = hwio.PinMode(ledPin, hwio.OUTPUT)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for {
		hwio.DigitalWrite(ledPin, hwio.HIGH)
		//hwio.Delay(100)
		hwio.DigitalWrite(ledPin, hwio.LOW)
		//hwio.Delay(100)
	}
}
