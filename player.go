package main

import (
	"time"

	pixgl "github.com/faiface/pixel/pixelgl"
)

func NewPlayer(win *pixgl.Window) *Entity {
	plr := NewEntity("sprites/player.png", win)

	plr.expands = append(plr.expands, &PlayerControl{PlayerSpeed})

	var plrMissilePool []*Entity
	for i := 0; i < PlayerMissilePoolSize; i++ {
		missile := NewEntity("sprites/missile.png", win)
		missile.active = false
		plrMissilePool = append(plrMissilePool, missile)
	}
	plrShooting := PlayerShooting{
		missileSpeed: PlayerMissileSpeed,
		pool:         plrMissilePool,
		cooldown:     PlayerCooldown,
	}
	plr.expands = append(plr.expands, &plrShooting)

	plr.pos.X = win.Bounds().W() / 2
	plr.pos.Y = win.Bounds().H() / 16

	return plr
}

type PlayerControl struct {
	speed float64
}

func (pc *PlayerControl) ExpandUpdate(p *Entity) error {
	if p.KeyPressed(pixgl.KeyA) {
		if !(p.pos.X-p.dim.width/2 < 0) {
			p.pos.X -= pc.speed
		}
	}
	if p.win.Pressed(pixgl.KeyD) {
		if !(p.pos.X+p.dim.width/2 > WindowWidth) {
			p.pos.X += pc.speed
		}
	}

	return nil
}

type PlayerShooting struct {
	missileSpeed float64
	pool         []*Entity
	cooldown     time.Duration
	lastShot     time.Time
}

func (ps *PlayerShooting) ExpandUpdate(p *Entity) error {
	for _, missile := range ps.pool {
		if missile.active {
			missile.pos.Y += ps.missileSpeed
			if missile.pos.Y > WindowHeight {
				missile.active = false
			} else {
				missile.Draw()
			}
		}
	}
	if p.KeyPressed(pixgl.KeySpace) || p.KeyPressed(pixgl.KeyS) {
		if time.Since(ps.lastShot) >= ps.cooldown {
			var missile *Entity

			// find un active missile and activate it after putting it in the right place
			for i, missileInPool := range ps.pool {
				if !missileInPool.active {
					missile = ps.pool[i]
					break
				}
			}

			missile.active = true
			missile.pos = p.pos

			ps.lastShot = time.Now()
		}
	}
	return nil
}
