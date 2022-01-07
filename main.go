package main

import (
	_ "image/png" // necessary for initialization

	pixgl "github.com/faiface/pixel/pixelgl"
)

func main() {
	pixgl.Run(run) // necessary for pixel to work something with threads idk
}

func run() {
	game := Initialize("Space shooters!")
	game.InitMenu()

	for !game.quit {
		game.quit = game.win.Closed()
		switch game.state {
		case StateMenu:
			game.Menu()
		case StatePlaySingal:
			game.SignalPlayerGame()
		case StateSignalOver:
			game.EndSingalPlayer()
		case StateMultiPlayer:
			game.state = StateMenu
		case StateQuit:
			game.quit = true
		}
	}
}
