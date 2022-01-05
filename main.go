package main

import (
	_ "image/png"

	"time"
	"math/rand"

	pixgl "github.com/faiface/pixel/pixelgl"
)

func main() {
	pixgl.Run(run) // necessary for pixel to work something with threads idk
}

func run() {
	game := Initialize("Space shooter")
	rand.Seed(time.Now().UnixNano())

	plr := NewPlayer(game)

	// annoyingThing := NewEntity("sprites/annoyingThing.png", game.win)
	// annoyingThing.pos = game.win.Bounds().Center()

	game.AddEntity(plr)
	// game.AddEntity(annoyingThing)

	for !game.win.Closed() {
		game.GameLoop()
		// touch := CheckCollision(plr, annoyingThing)
		// if touch {
		// 	fmt.Printf("touched?: %v\n", touch)
		// }
	}
}
