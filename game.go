package main

import (
	"image/color"
	"time"

	pixgl "github.com/faiface/pixel/pixelgl"
)

type Game struct {
	wcfg     *pixgl.WindowConfig
	win      *pixgl.Window
	framer   *FrameLimiter
	entities []*Entity
}

func (g *Game) GameLoop() {
	g.win.Clear(color.RGBA{30, 30, 30, 255})

	for _, e := range g.entities {
		e.Update()
		e.Draw()
	}

	g.framer.SetTitleWithFPS(g.win, g.wcfg)

	g.win.Update()

	time.Sleep(time.Millisecond * MillersecondsPerFrame)
}

func (g *Game) AddEntity(e *Entity) {
	g.entities = append(g.entities, e)
}
