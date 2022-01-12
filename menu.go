package main

import (
	"image/color"

	pix "github.com/faiface/pixel"
	pixgl "github.com/faiface/pixel/pixelgl"
)

const (
	UIButtonFilePath      = "menu/button.png"
	UIButtonFocusFilePath = "menu/focusButton.png"
	UIButtonWidth         = 360
	UIButtonHeight        = 80

	BackButtonFilePath      = "menu/backButton.png"
	BackButtonFocusFilePath = "menu/focusBackButton.png"
)

func NewBackButton(g *Game, pos pix.Vec) *Entity {
	e := NewEntity(BackButtonFilePath, g.win)
	e.pos = pos

	bb := &BackButton{}
	bb.g = g
	bb.fb = NewEntity(BackButtonFocusFilePath, g.win)

	e.expands = append(e.expands, bb)

	return e
}

type BackButton struct {
	g             *Game
	fb            *Entity
	shouldTrigger bool
}

func (bb *BackButton) ExpandUpdate(e *Entity) error {
	if bb.shouldTrigger && bb.g.win.Pressed(pixgl.MouseButtonLeft) {
		bb.g.ChangeState(StateMenu)
	}

	mousePos := e.win.MousePosition()
	if EToR(e, 0).Contains(mousePos) {
		bb.fb.pos = e.pos
		bb.fb.Draw()
		if bb.g.win.Pressed(pixgl.MouseButtonLeft) {
			bb.shouldTrigger = true
		}
	}

	return nil
}

func NewUIButton(g *Game, toState int, label string, pos pix.Vec) *Entity {
	e := NewEntity(UIButtonFilePath, g.win)
	e.pos = pos

	uib := &UIButton{}
	uib.g = g
	uib.fb = NewEntity(UIButtonFocusFilePath, g.win)
	uib.toState = toState
	uib.label = label
	uib.shouldTrigger = false

	e.expands = append(e.expands, uib)
	return e
}

type UIButton struct {
	g             *Game
	fb            *Entity
	toState       int
	label         string
	shouldTrigger bool
}

func (uib *UIButton) ExpandUpdate(e *Entity) error {
	if uib.shouldTrigger && !uib.g.win.Pressed(pixgl.MouseButtonLeft) {
		uib.g.ChangeState(uib.toState)
	}

	mousePos := e.win.MousePosition()
	if EToR(e, 0).Contains(mousePos) {
		uib.fb.pos = e.pos
		uib.fb.Draw()
		if uib.g.win.Pressed(pixgl.MouseButtonLeft) {
			uib.shouldTrigger = true
		}
	}

	uib.g.tw.WriteText(uib.label, pix.V(e.pos.X-UIButtonWidth/2.5, e.pos.Y-UIButtonHeight/5), color.RGBA{60, 60, 60, 255}, 1)

	return nil
}

func InitMenu(g *Game) {
	bmh := 2.5

	g.entities = []*Entity{}

	logo := NewEntity("menu/logo.png", g.win)
	logo.pos = pix.V(WindowWidth/2, WindowHeight-222)

	b1 := NewUIButton(g, StateSingalPlayer, "Single Player", pix.V(WindowWidth/2, WindowHeight/bmh))
	// b2 := NewUIButton(g, StateHighScore, "High Score", pix.V(WindowWidth/2, WindowHeight/bmh-85))
	b3 := NewUIButton(g, StateMultiPlayer, "Multiplayer", pix.V(WindowWidth/2, WindowHeight/bmh-85))
	b4 := NewUIButton(g, StateQuit, "Quit", pix.V(WindowWidth/2, WindowHeight/bmh-170))

	g.AddEntity(logo)
	g.AddEntity(b1)
	// g.AddEntity(b2)
	g.AddEntity(b3)
	g.AddEntity(b4)
}

func (g *Game) Menu() {
	g.WinClear()
	for _, e := range g.entities {
		e.Draw()
		e.Update()
	}
	g.WinUpdate()
}

func InitQuitState(g *Game) {
	g.quit = true
}
