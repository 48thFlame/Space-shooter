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

func InitSingalPlayerOverState(g *Game) {
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
