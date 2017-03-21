package main

import (
	. "github.com/alexellis/rpi"
)

func main() {
	WiringPiSetup()

	//use default pin naming
	PinMode(PIN_GPIO_7, OUTPUT)
  for{
	  DigitalWrite(PIN_GPIO_7, LOW)
	  //Delay(400)
	  DigitalWrite(PIN_GPIO_7, HIGH)
  }

	//use raspberry pi board pin numbering, similiar to RPi.GPIO.setmode(RPi.GPIO.BOARD)
	Delay(400)
	DigitalWrite(BoardToPin(7), LOW)
	Delay(400)
	DigitalWrite(BoardToPin(7), HIGH)

	//use raspberry pi bcm gpio numbering, similiar to RPi.GPIO.setmode(RPi.GPIO.BCM)
	Delay(400)
	DigitalWrite(GpioToPin(4), LOW)
	Delay(400)
	DigitalWrite(GpioToPin(4), HIGH)

	Delay(400)
	DigitalWrite(PIN_GPIO_7, LOW)
}
