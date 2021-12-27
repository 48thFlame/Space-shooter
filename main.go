package main

import (
	"fmt"
	"image/color"
	_ "image/png"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

func main() {
	pixelgl.Run(run)
}

func run() {
	wcfg, win, framer := Initialize("Space shooter!")

	playerImg, err := LoadPicture("sprites/player.png")
	if err != nil {
		panic(err)
	}
	playerSprite := pixel.NewSprite(playerImg, playerImg.Bounds())
	angle := 0.0

	var (
		frames = 0
		second = time.Tick(time.Second)
	)

	for !win.Closed() {
		if framer.ShouldDoNextFrame() {
			win.Clear(color.RGBA{30, 30, 30, 255})

			angle += 0.05

			mat := pixel.IM
			mat = mat.Moved(win.Bounds().Center())
			mat = mat.Rotated(win.Bounds().Center(), angle)

			playerSprite.Draw(win, mat)

			win.Update()

			frames++
			select {
			case <-second:
				win.SetTitle(fmt.Sprintf("%s | FPS: %d", wcfg.Title, frames))
				frames = 0
			default:
			}
		}
	}

}
