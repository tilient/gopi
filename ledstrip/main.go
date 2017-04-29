package main

import (
	"github.com/tilient/gopi/ws2811"
	"time"
)

// ===========================================================
// Clock
// ===========================================================

const numPixels = 15

func color(r, g, b uint32) uint32 {
	return (g << 16) + (r << 8) + b
}

func wheel(pos uint32) uint32 {
	// Generate rainbow colors across 0-255 positions.
	if pos < 85 {
		return color(pos*3, 255-pos*3, 0)
	}
	pos -= 85
	if pos < 85 {
		return color(255-pos*3, 0, pos*3)
	}
	pos -= 85
	return color(0, pos*3, 255-pos*3)
}

func rainbow() {
	// Draw rainbow that fades across all pixels at once.
	for j := 0; j < 256; j++ {
		for i := 0; i < numPixels; i++ {
			ws2811.SetLed(i, wheel(uint32((i*16+j)%255)))
		}
		ws2811.Render()
		ws2811.Wait()
		time.Sleep(10 * time.Millisecond)
	}
}

func main() {
	ws2811.Init(18, numPixels, 55)
	defer ws2811.Fini()

	ws2811.Clear()
	ws2811.Render()
	ws2811.Wait()
	for {
		rainbow()
	}
}

// ===========================================================
