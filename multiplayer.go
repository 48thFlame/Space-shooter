package main

import (
	"fmt"
	"image/color"

	pix "github.com/faiface/pixel"
)

func InitPlrHealths(g *Game) {
	newMultiHealths := make(map[int]int)
	for k, v := range MultiHealths {
		newMultiHealths[k] = v
	}
	g.plrHealths = newMultiHealths
}

func InitMultiPlayerState(g *Game) {
	g.entities = []*Entity{}
	InitPlrHealths(g)

	backButton := NewBackButton(g, pix.V(52, WindowHeight-52))

	plr1, ps1 := NewMultiPlayer(g, 1)
	plr2, ps2 := NewMultiPlayer(g, 2)

	AddMPMCToPlr(plr1, ps2, ps1, g, 1, 2)
	AddMPMCToPlr(plr2, ps1, ps2, g, 2, 1)

	g.AddEntity(backButton)

	g.AddEntity(plr1)
	g.AddEntity(plr2)
}

func (g *Game) MultiPlayer() {
	g.WinClear()

	g.tw.WriteText(fmt.Sprintf("%v", g.plrHealths[1]), pix.V(WindowWidth/2-32, WindowHeight-64), color.RGBA{52, 152, 219, 255}, 1)
	g.tw.WriteText(fmt.Sprintf("%v", g.plrHealths[2]), pix.V(WindowWidth/2+32, WindowHeight-64), color.RGBA{155, 89, 182, 255}, 1)

	for _, e := range g.entities {
		e.Update()
		e.Draw()
	}

	g.framer.SetTitleWithFPS(g.win, g.wcfg)

	g.WinUpdate()

}

func InitMultiPlayerOverState(g *Game) {
	g.entities = []*Entity{}
	InitPlrHealths(g)

	gameOver := NewEntity("menu/gameOver.png", g.win)
	gameOver.pos = pix.V(WindowWidth/2, WindowHeight-WindowHeight/4)

	button := NewUIButton(g, StateMenu, "Menu", pix.V(WindowWidth/2, WindowHeight/2-300))

	g.AddEntity(gameOver)
	g.AddEntity(button)
}

func (g *Game) MultiPlayerOver() {
	g.WinClear()
	for _, e := range g.entities {
		e.Draw()
		e.Update()
	}

	g.tw.WriteText(
		fmt.Sprintf("The %v player won!", g.multiWinnerColor),
		pix.V(336, 453),
		color.RGBA{255, 255, 255, 255},
		1,
	)

	g.WinUpdate()
}
