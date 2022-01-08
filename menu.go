package main

import (
	pix "github.com/faiface/pixel"
	pixgl "github.com/faiface/pixel/pixelgl"
)

type MouseFollower struct{}

func (mf *MouseFollower) ExpandUpdate(e *Entity) error {
	e.pos = e.win.MousePosition()
	return nil
}

func NewButton(filePath string, pos pix.Vec, g *Game, id int) *Entity {
	e := NewEntity(filePath, g.win)

	e.pos = pos
	mb := &UIButton{}
	mb.id = id
	mb.g = g

	e.expands = append(e.expands, mb)

	return e
}

type UIButton struct {
	id    int
	g     *Game
	label string
}

func (mb *UIButton) ExpandUpdate(e *Entity) error {
	if mb.g.win.Pressed(pixgl.MouseButtonLeft) {
		mousePos := e.win.MousePosition()
		if EToR(e, 0).Contains(mousePos) {
			mb.g.ChangeState(mb.id)
		}
	}
	return nil
}

func InitMenu(g *Game) {
	bmh := 2.75

	g.entities = []*Entity{}
	background := NewEntity("menu/menu.png", g.win)
	background.pos = pix.V(WindowWidth/2, WindowHeight/2)

	mouse := NewEntity("menu/mouse.png", g.win)
	mouse.expands = append(mouse.expands, &MouseFollower{})

	b1 := NewButton("menu/button 1.png", pix.V(WindowWidth/2, WindowHeight/bmh), g, StateSingalPlayer)
	b2 := NewButton("menu/button 2.png", pix.V(WindowWidth/2, WindowHeight/bmh-56), g, StateMultiPlayer)
	b3 := NewButton("menu/button 3.png", pix.V(WindowWidth/2, WindowHeight/bmh-112), g, StateQuit)

	g.AddEntity(background)

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

func InitQuitState(g *Game) {
	g.quit = true
}
