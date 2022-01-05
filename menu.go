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
	mb := &MenuButton{}
	mb.id = id
	mb.g = g

	e.expands = append(e.expands, mb)

	return e
}

type MenuButton struct {
	id int
	g  *Game
}

func (mb *MenuButton) ExpandUpdate(e *Entity) error {
	if mb.g.win.Pressed(pixgl.MouseButtonLeft) {
		mousePos := e.win.MousePosition()
		if EToR(e).Contains(mousePos) {
			if mb.id == StatePlaySingal {
				mb.g.InitSignalPlayerGame()
			} else if mb.id == StateMenu {
				mb.g.InitMenu()
			}
			mb.g.state = mb.id
		}
	}
	return nil
}
