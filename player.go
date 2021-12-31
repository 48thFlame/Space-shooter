package main

import (
	"time"

	pixgl "github.com/faiface/pixel/pixelgl"
)

const (
	PlayerAccSpeed        = 0.93
	// PlayerMaxSpeed        = 8
	PlayerSlowSpeed       = 0.11
	PlayerDragSpeed       = 3 // how much of player speed should we slow down
	PlayerRotSpeed        = 0.07
	PlayerCooldown        = time.Millisecond * 125
	PlayerMissileSpeed    = 12
	PlayerMissilePoolSize = 80
)

func NewPlayer(win *pixgl.Window) *Entity {
	plr := NewEntity("sprites/player.png", win)

	plr.expands = append(plr.expands, &PlayerControl{speed: PlayerAccSpeed})

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
	plr.pos.Y = win.Bounds().H() / 5

	return plr
}

type PlayerControl struct {
	speed  float64
	xv, yv float64
}

func (pc *PlayerControl) ExpandUpdate(p *Entity) error {
	if KeyPressed(p, pixgl.KeyA) {
		p.rot += PlayerRotSpeed
	}
	if KeyPressed(p, pixgl.KeyD) {
		p.rot -= PlayerRotSpeed
	}

	if KeyPressed(p, pixgl.KeyW) {
		sxv, syv := ChangeXY(p.rot)
		pc.xv += sxv * PlayerAccSpeed
		pc.yv += syv * PlayerAccSpeed
	}

	// change x, y for player
	if p.pos.X > WindowWidth+p.dim.width {
		p.pos.X = 0
	} else if p.pos.X < 0-p.dim.width {
		p.pos.X = WindowWidth
	}
	p.pos.X += pc.xv

	if p.pos.Y > WindowHeight+p.dim.height {
		p.pos.Y = 0
	} else if p.pos.Y < 0-p.dim.height {
		p.pos.Y = WindowHeight
	}
	p.pos.Y += pc.yv

	// slow down player
	if pc.xv > 0 {
		pc.xv -= PlayerSlowSpeed * pc.xv / PlayerDragSpeed
	}
	if pc.xv < 0 {
		pc.xv += PlayerSlowSpeed * pc.xv / -PlayerDragSpeed
	}

	if pc.yv > 0 {
		pc.yv -= PlayerSlowSpeed * pc.yv / PlayerDragSpeed
	}
	if pc.yv < 0 {
		pc.yv += PlayerSlowSpeed * pc.yv / -PlayerDragSpeed
	}

	return nil
}

type PlayerShooting struct {
	missileSpeed float64
	pool         []*Entity
	cooldown     time.Duration
	lastShot     time.Time
}

func (ps *PlayerShooting) MoveMissileOffPlayer(missile, plr *Entity) {
	xv, yv := ChangeXY(missile.rot)
	missile.pos.X += xv * plr.dim.width / 2
	missile.pos.Y += yv * plr.dim.height / 2
}

func (ps *PlayerShooting) ExpandUpdate(p *Entity) error {
	for _, missile := range ps.pool {
		if missile.active {
			// missile.pos.Y += ps.missileSpeed
			xv, yv := ChangeXY(missile.rot)
			missile.pos.X += xv * ps.missileSpeed
			missile.pos.Y += yv * ps.missileSpeed
			if missile.pos.Y > WindowHeight || missile.pos.Y < 0 || missile.pos.X > WindowWidth || missile.pos.X < 0{
				missile.active = false
				} else {
					missile.Draw()
				}
			}
		}
		if KeyPressed(p, pixgl.KeySpace) || KeyPressed(p, pixgl.KeyS) {
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
				missile.rot = p.rot
				missile.pos = p.pos
				ps.MoveMissileOffPlayer(missile, p)
		
				ps.lastShot = time.Now()
			}
		}
		return nil
	}
