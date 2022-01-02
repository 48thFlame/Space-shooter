package main

import (
	"image/color"

	pixgl "github.com/faiface/pixel/pixelgl"
)

type Game struct {
	wcfg     *pixgl.WindowConfig
	win      *pixgl.Window
	framer   *FrameLimiter
	entities []*Entity
}

func (g *Game) GameLoop() {
	g.framer.InitFrameCounter()

	if g.framer.ShouldDoNextFrame() {
		g.win.Clear(color.RGBA{30, 30, 30, 255})

		for _, e := range g.entities {
			e.Update()
			e.Draw()
		}

		g.framer.SetTitleWithFPS(g.win, g.wcfg)

		g.win.Update()
	}
}

func (g *Game) AddEntity(e *Entity) {
	g.entities = append(g.entities, e)
	// if collision {
	// 	g.collisionEntities = append(g.collisionEntities, e)
	// }
}
