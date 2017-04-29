package main

import (
	"github.com/tilient/gopi/ws2811"
)

// ===========================================================
// Clock
// ===========================================================

func color(r, g, b uint32) uint32 {
	return (g << 16) + (r << 8) + b
}

func main() {
	ws2811.Init(18, 15, 255)
	defer ws2811.Fini()

	ws2811.Clear()
	ws2811.Render()
	ws2811.Wait()
	ws2811.SetLed(3, color(100, 0, 0))
	ws2811.SetLed(2, color(0, 100, 0))
	ws2811.SetLed(1, color(0, 0, 100))

	ws2811.SetLed(7, color(100, 100, 0))
	ws2811.SetLed(6, color(100, 0, 100))
	ws2811.SetLed(5, color(0, 100, 100))
	ws2811.Render()
	ws2811.Wait()
}

// ===========================================================
