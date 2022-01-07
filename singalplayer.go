package main

import (
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
