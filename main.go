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
	game.ChangeState(StateMenu)

	for !game.quit {
		game.quit = game.win.Closed()
		switch game.state {
		case StateMenu:
			game.Menu()
		case StateSingalPlayer:
			game.SignalPlayerGame()
		case StateSingalPlayerOver:
			game.EndSingalPlayer()
		case StateMultiPlayer:
			game.state = StateMenu
		}
	}
}
