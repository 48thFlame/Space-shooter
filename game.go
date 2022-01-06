package main

import (
	"fmt"
	"image/color"
	"time"

	pix "github.com/faiface/pixel"
	pixgl "github.com/faiface/pixel/pixelgl"
)

type Game struct {
	quit     bool
	state    int
	wcfg     *pixgl.WindowConfig
	win      *pixgl.Window
	framer   *FrameCounter
	entities []*Entity
	tw       *TextWriter
	score    int
	level    int
}

func (g *Game) WinClear() {
	g.win.Clear(color.RGBA{30, 30, 30, 255})
}

func (g *Game) WinUpdate() {
	g.win.Update()
	time.Sleep(time.Millisecond * MillersecondsPerFrame)
}

func (g *Game) InitSignalPlayerGame() {
	g.entities = []*Entity{}
	g.score = 0
	g.level = 1
	plr := NewPlayer(g)

	g.AddEntity(plr)
}

func (g *Game) SignalPlayerGame() {
	g.WinClear()

	g.tw.WriteText(fmt.Sprint(g.score), pix.V(WindowWidth-128, WindowHeight-82), color.RGBA{220, 220, 220, 255})

	for _, e := range g.entities {
		e.Update()
		e.Draw()
	}

	g.framer.SetTitleWithFPS(g.win, g.wcfg)

	g.WinUpdate()
}

func (g *Game) InitEndSingalPlayer() {
	g.entities = []*Entity{}
	menu := NewEntity("menu/game over.png", g.win)
	menu.pos = pix.V(WindowWidth/2, WindowHeight/2)

	button := NewButton("menu/button 4.png", pix.V(WindowWidth/2, WindowHeight/2-300), g, StateMenu)

	g.AddEntity(menu)

	g.AddEntity(button)
}

func (g *Game) EndSingalPlayer() {
	g.WinClear()
	for _, e := range g.entities {
		e.Update()
		e.Draw()
	}

	g.tw.WriteText(
		fmt.Sprint(g.score),
		pix.V(750, 453),
		color.RGBA{255, 255, 255, 255},
	)

	g.WinUpdate()
}

func (g *Game) InitMenu() {
	bmh := 2.75

	g.entities = []*Entity{}
	menu := NewEntity("menu/menu.png", g.win)
	menu.pos = pix.V(WindowWidth/2, WindowHeight/2)

	mouse := NewEntity("menu/mouse.png", g.win)
	mouse.expands = append(mouse.expands, &MouseFollower{})

	b1 := NewButton("menu/button 1.png", pix.V(WindowWidth/2, WindowHeight/bmh), g, StatePlaySingal)
	b2 := NewButton("menu/button 2.png", pix.V(WindowWidth/2, WindowHeight/bmh-56), g, StateMultiPlayer)
	b3 := NewButton("menu/button 3.png", pix.V(WindowWidth/2, WindowHeight/bmh-112), g, StateQuit)

	g.AddEntity(menu)

	g.AddEntity(b1)
	g.AddEntity(b2)
	g.AddEntity(b3)

	g.AddEntity(mouse)
}

func (g *Game) Menu() {
	g.WinClear()
	for _, e := range g.entities {
		e.Update()
		e.Draw()
	}
	g.WinUpdate()
}

func (g *Game) AddEntity(e *Entity) {
	g.entities = append(g.entities, e)
}
