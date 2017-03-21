package main

import (
	"flag"
	"image/color"
	"time"

	"github.com/mcuadros/go-rpi-ws281x"
)

var (
	pin        = flag.Int("gpio-pin", 18, "GPIO pin")
	width      = flag.Int("width", 4, "LED matrix width")
	height     = flag.Int("height", 2, "LED matrix height")
	brightness = flag.Int("brightness", 64, "Brightness (0-255)")
)

func main() {
	config := ws281x.DefaultConfig
	config.Brightness = *brightness
	config.Pin = *pin

	c, err := ws281x.NewCanvas(*width, *height, &config)
	fatal(err)

	defer c.Close()
	err = c.Initialize()
	fatal(err)

	bounds := c.Bounds()

	for r := 0; r < 256; r++ {
		for g := 0; g < 128; g += 8 {
			color := color.RGBA{uint8(r), uint8(g), 0, 255}
			for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
				for x := bounds.Min.X; x < bounds.Max.X; x++ {
					c.Set(x, y, color)
					c.Render()
				}
			}
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func init() {
	flag.Parse()
}

func fatal(err error) {
	if err != nil {
		panic(err)
	}
}
