package main

import (
	_ "image/png"

	"github.com/faiface/pixel/pixelgl"
)

func main() {
	pixelgl.Run(run) // necessary for pixel to work something with threads idk
}

func run() {
	game := Initialize("Space shooter")

	plr := NewPlayer(game.win)
	game.entities = append(game.entities, plr)

	game.GameLoop()
}
