package main

import (
	"math"
	"time"

	pixgl "github.com/faiface/pixel/pixelgl"
)

const (
	PlayerAccSpeed                 = 0.66
	PlayerMaxSpeed                 = 16
	PlayerPercentOfSpeedShouldSlow = 0.02
	PlayerRotSpeed                 = 0.11
	PlayerCooldown                 = time.Millisecond * 125
	PlayerMissileSpeed             = 25
	PlayerMissilePoolSize          = 80
)

func NewPlayer(win *pixgl.Window) *Entity {
	plr := NewEntity("sprites/player.png", win)

	plrControl := PlayerControl{
		speed:   PlayerAccSpeed,
		forward: pixgl.KeyW,
		left:    pixgl.KeyA,
		right:   pixgl.KeyD,
	}
	plr.expands = append(plr.expands, &plrControl)

	var plrMissilePool []*Entity
	for i := 0; i < PlayerMissilePoolSize; i++ {
		missile := NewEntity("sprites/missile.png", win)
		missile.active = false
		plrMissilePool = append(plrMissilePool, missile)
	}
	plrShooting := PlayerShooting{
		shoot:        pixgl.KeyS,
		missileSpeed: PlayerMissileSpeed,
		pool:         plrMissilePool,
		cooldown:     PlayerCooldown,
	}
	plr.expands = append(plr.expands, &plrShooting)

	fire := NewEntity("sprites/fire.png", win)
	plrFire := PlayerFire{
		pc:   &plrControl,
		fire: fire,
	}
	plr.expands = append(plr.expands, &plrFire)

	// plr.pos = win.Bounds().Center()

	return plr
}

type PlayerControl struct {
	speed  float64
	xv, yv float64
	// keys:
	forward, left, right pixgl.Button
}

func (pc *PlayerControl) ExpandUpdate(p *Entity) error {
	if KeyPressed(p, pc.left) {
		p.rot += PlayerRotSpeed
	}
	if KeyPressed(p, pc.right) {
		p.rot -= PlayerRotSpeed
	}

	if KeyPressed(p, pc.forward) {
		sxv, syv := ChangeXY(p.rot)
		if pc.xv < PlayerMaxSpeed && pc.xv > -PlayerMaxSpeed {
			pc.xv += sxv * PlayerAccSpeed
		}
		if pc.yv < PlayerMaxSpeed && pc.yv > -PlayerMaxSpeed {
			pc.yv += syv * PlayerAccSpeed
		}
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
		pc.xv -= PlayerPercentOfSpeedShouldSlow * pc.xv
	}
	if pc.xv < 0 {
		pc.xv += PlayerPercentOfSpeedShouldSlow * -pc.xv
	}

	if pc.yv > 0 {
		pc.yv -= PlayerPercentOfSpeedShouldSlow * pc.yv
	}
	if pc.yv < 0 {
		pc.yv += PlayerPercentOfSpeedShouldSlow * -pc.yv
	}

	return nil
}

type PlayerShooting struct {
	shoot        pixgl.Button
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
			if missile.pos.Y > WindowHeight || missile.pos.Y < 0 || missile.pos.X > WindowWidth || missile.pos.X < 0 {
				missile.active = false
			} else {
				missile.Draw()
			}
		}
	}
	if KeyPressed(p, ps.shoot) {
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

type PlayerFire struct {
	pc   *PlayerControl
	fire *Entity
}

func (pf *PlayerFire) MoveToFireBackOfPlayer(fire, plr *Entity) {
	xv, yv := ChangeXY(fire.rot)
	fire.pos = plr.pos
	fire.pos.X += xv * plr.dim.width / 2
	fire.pos.Y += yv * plr.dim.height / 2
	// fire.pos = pixel.V(100, 100)
}

func (pf *PlayerFire) ExpandUpdate(p *Entity) error {
	// fmt.Println("here")
	pf.fire.rot = p.rot + math.Pi
	pf.MoveToFireBackOfPlayer(pf.fire, p)

	if KeyPressed(p, pf.pc.forward) {
		pf.fire.active = true

	} else {
		pf.fire.active = false
	}

	pf.fire.Draw()

	return nil
}
