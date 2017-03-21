package main

import (
	//"time"

	"github.com/kidoman/embd"
	_ "github.com/kidoman/embd/host/rpi"
)

func main() {
  embd.InitGPIO()
  defer embd.CloseGPIO()

  pin, _ := embd.NewDigitalPin(4)

  pin.SetDirection(embd.Out)
  for {
    pin.Write(embd.High)
    //time.Sleep(100 * time.Millisecond)
    pin.Write(embd.Low)
    //time.Sleep(100 * time.Millisecond)
  }
}
