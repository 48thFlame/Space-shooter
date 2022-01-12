package main

import (
	"fmt"
	"image/color"

	pix "github.com/faiface/pixel"
	pixgl "github.com/faiface/pixel/pixelgl"
)

func NewSingalPlayer(g *Game) *Entity {
	plr := NewEntity("sprites/player1.png", g.win)

	plrControl := PlayerControl{
		speed:   PlayerAccSpeed,
		forward: pixgl.KeyW,
		left:    pixgl.KeyA,
		right:   pixgl.KeyD,
	}
	plr.expands = append(plr.expands, &plrControl)

	var plrMissilePool []*Entity
	for i := 0; i < PlayerMissilePoolSize; i++ {
		missile := NewEntity("sprites/missile.png", g.win)
		plrMissilePool = append(plrMissilePool, missile)
	}
	plrShooting := &PlayerShooting{
		shoot:        pixgl.KeyS,
		missileSpeed: PlayerMissileSpeed,
		pool:         plrMissilePool,
		cooldown:     PlayerCooldown,
	}
	plr.expands = append(plr.expands, plrShooting)

	fire := NewEntity("sprites/fire.png", g.win)
	plrFire := PlayerFire{
		pc:   &plrControl,
		fire: fire,
	}
	plr.expands = append(plr.expands, &plrFire)

	prc := NewRockController(plr, g)
	plr.expands = append(plr.expands, prc)

	pmrc := &PlayerMissileRockCollision{}
	pmrc.ps = plrShooting
	pmrc.prc = prc
	pmrc.g = g
	pmrc.plr = plr
	plr.expands = append(plr.expands, pmrc)

	plr.pos = g.win.Bounds().Center()

	return plr
}

func InitSignalPlayerState(g *Game) {
	g.entities = []*Entity{}
	g.score = 0
	g.level = 1
	plr := NewSingalPlayer(g)
	bb := NewBackButton(g, pix.V(56, WindowHeight-56))

	g.AddEntity(bb)
	g.AddEntity(plr)
}

func (g *Game) SignalPlayerGame() {
	g.WinClear()

	g.tw.WriteText(fmt.Sprint(g.score), pix.V(WindowWidth-128, WindowHeight-64), color.RGBA{220, 220, 220, 255}, 1)

	for _, e := range g.entities {
		e.Update()
		e.Draw()
	}

	g.framer.SetTitleWithFPS(g.win, g.wcfg)

	g.WinUpdate()
}

func InitSingalPlayerOverState(g *Game) {
	g.entities = []*Entity{}
	gameOver := NewEntity("menu/gameOver.png", g.win)
	gameOver.pos = pix.V(WindowWidth/2, WindowHeight-WindowHeight/4)

	button := NewUIButton(g, StateMenu, "Menu", pix.V(WindowWidth/2, WindowHeight/2-300))

	g.AddEntity(gameOver)

	g.AddEntity(button)
}

func (g *Game) EndSingalPlayer() {
	g.WinClear()
	for _, e := range g.entities {
		e.Draw()
		e.Update()
	}

	g.tw.WriteText(
		fmt.Sprintf("Your score was: %v", g.score),
		pix.V(336, 453),
		color.RGBA{255, 255, 255, 255},
		1,
	)

	g.WinUpdate()
}

// func InitHighScoreState(g *Game) {
// 	g.entities = []*Entity{}
// }

// func (g *Game) HighScore() {
// 	g.ChangeState(StateMenu)
// }
