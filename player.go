package main

import (
	"math"
	"time"

	pixgl "github.com/faiface/pixel/pixelgl"
)

const (
	PlayerAccSpeed                      = 0.66
	PlayerMaxSpeed                      = 16
	PlayerPercentOfSpeedShouldSlow      = 0.02
	PlayerRotSpeed                      = 0.17
	PlayerCooldown                      = time.Millisecond * 250
	PlayerMissileSpeed                  = 25
	PlayerMissilePoolSize               = 50	
)

// func NewPlayer(filePath string, g *Game) *Entity {
	
// }

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
	activePool   []*Entity
	cooldown     time.Duration
	lastShot     time.Time
}

func (ps *PlayerShooting) MoveMissileFrontOfPlayer(missile, plr *Entity) {
	xv, yv := ChangeXY(missile.rot)
	missile.pos.X += xv * plr.dim.width / 2
	missile.pos.Y += yv * plr.dim.height / 2
}

func (ps *PlayerShooting) ExpandUpdate(p *Entity) error {
	var IToBeMoved []int

	for missileI, missile := range ps.activePool {
		xv, yv := ChangeXY(missile.rot)
		missile.pos.X += xv * ps.missileSpeed
		missile.pos.Y += yv * ps.missileSpeed

		if missile.pos.Y > WindowHeight || missile.pos.Y < 0 || missile.pos.X > WindowWidth || missile.pos.X < 0 {
			IToBeMoved = append(IToBeMoved, missileI)
		} else {
			missile.Draw()
		}
	}

	var newActivePool []*Entity
	for i := 0; i < len(ps.activePool); i++ {
		if Contains(IToBeMoved, i) {
			ps.pool = append(ps.pool, ps.activePool[i])
		} else {
			newActivePool = append(newActivePool, ps.activePool[i])
		}
	}

	ps.activePool = newActivePool

	if KeyPressed(p, ps.shoot) {
		if time.Since(ps.lastShot) >= ps.cooldown {
			if len(ps.pool) > 0 {
				missile := ps.pool[0]
				ps.pool = RemoveFromESlice(ps.pool, 0)

				missile.rot = p.rot
				missile.pos = p.pos
				ps.MoveMissileFrontOfPlayer(missile, p)

				ps.activePool = append(ps.activePool, missile)
				ps.lastShot = time.Now()
			}
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
}

func (pf *PlayerFire) ExpandUpdate(p *Entity) error {
	pf.fire.rot = p.rot + math.Pi
	pf.MoveToFireBackOfPlayer(pf.fire, p)

	if KeyPressed(p, pf.pc.forward) {
		pf.fire.Draw()
	}

	return nil
}
