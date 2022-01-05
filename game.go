package main

import (
	"fmt"
	"image/color"
	"time"

	pix "github.com/faiface/pixel"
	pixgl "github.com/faiface/pixel/pixelgl"
)

type Game struct {
	wcfg     *pixgl.WindowConfig
	win      *pixgl.Window
	framer   *FrameCounter
	entities []*Entity
	tw       *TextWriter
	score    int
	level    int
}

func (g *Game) GameLoop() {
	g.win.Clear(color.RGBA{30, 30, 30, 255})

	g.tw.WriteText(fmt.Sprint(g.score), pix.V(WindowWidth-128, WindowHeight-82), color.RGBA{220, 220, 220, 255})

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
