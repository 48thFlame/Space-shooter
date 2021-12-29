package main

import (
	"image/color"
	_ "image/png"

	"github.com/faiface/pixel/pixelgl"
)

func main() {
	pixelgl.Run(run) // necessary for pixel to work something with threads idk
}

func run() {
	wcfg, win, framer := Initialize("Space shooter!")
	framer.InitFrameCounter()

	plr := NewPlayer(win)

	for !win.Closed() {
		if framer.ShouldDoNextFrame() {
			win.Clear(color.RGBA{30, 30, 30, 255})

			plr.Update()
			plr.Draw()

			framer.SetTitleWithFPS(win, wcfg)
			
			win.Update()
		}
	}
}
