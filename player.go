package main

import (
	"fmt"
	"math"
	"time"

	pix "github.com/faiface/pixel"
	pixgl "github.com/faiface/pixel/pixelgl"
)

var (
	MultiKeyMap = map[int]map[rune]pixgl.Button{
		1: {
			'f': pixgl.KeyW,
			'l': pixgl.KeyA,
			'r': pixgl.KeyD,
			's': pixgl.KeyS,
		},
		2: {
			'f': pixgl.KeyUp,
			'l': pixgl.KeyLeft,
			'r': pixgl.KeyRight,
			's': pixgl.KeyDown,
		},
	}
	MultiInitPosMap = map[int]pix.Vec{
		1: pix.V(128, 128),
		2: pix.V(WindowWidth-128, WindowHeight-128),
	}
	MultiInitRotMap = map[int]float64{
		1: 1.75 * math.Pi,
		2: 0.75 * math.Pi,
	}
	MultiHealths = map[int]int{
		1: MultiPlayerStartHealth,
		2: MultiPlayerStartHealth,
	}
	MultiPlayerColors = map[int]string{
		1: "blue",
		2: "purple",
	}
)

const (
	PlayerAccSpeed                 = 0.475
	PlayerMaxSpeed                 = 12
	PlayerPercentOfSpeedShouldSlow = 0.01
	PlayerRotSpeed                 = 0.12
	PlayerCooldown                 = time.Millisecond * 250
	PlayerMissileSpeed             = 32
	PlayerMissilePoolSize          = 50
	MultiPlayerStartHealth         = 3
)

func NewMultiPlayer(g *Game, plrNum int) (*Entity, *PlayerShooting) {
	plr := NewEntity(fmt.Sprintf("sprites/player%v.png", plrNum), g.win)
	plr.pos = MultiInitPosMap[plrNum]
	plr.rot = MultiInitRotMap[plrNum]

	pc := &PlayerControl{}
	pc.speed = PlayerAccSpeed
	pc.forward = MultiKeyMap[plrNum]['f']
	pc.left = MultiKeyMap[plrNum]['l']
	pc.right = MultiKeyMap[plrNum]['r']

	pf := &PlayerFire{}
	pf.pc = pc
	pf.fire = NewEntity("sprites/fire.png", g.win)

	var plrMissilePool []*Entity
	for i := 0; i < PlayerMissilePoolSize; i++ {
		missile := NewEntity("sprites/missile.png", g.win)
		plrMissilePool = append(plrMissilePool, missile)
	}
	ps := &PlayerShooting{
		shoot:        MultiKeyMap[plrNum]['s'],
		missileSpeed: PlayerMissileSpeed,
		pool:         plrMissilePool,
		cooldown:     PlayerCooldown,
	}

	plr.expands = append(plr.expands, pc, pf, ps)

	return plr, ps
}

func MoveMultiPlayersToStart(plr1, plr2 *Entity) {
	plr1.pos = MultiInitPosMap[1]
	plr1.rot = MultiInitRotMap[1]

	plr2.pos = MultiInitPosMap[2]
	plr2.rot = MultiInitRotMap[2]
}

func AddMPMCToPlr(plr1 *Entity, otherPS, myPS *PlayerShooting, g *Game, plrNum, otherPlrNum int) {
	mpmc := &MultiPlayerMissileCollision{}
	mpmc.g = g
	mpmc.thisPlr = plr1
	mpmc.myPS = myPS
	mpmc.otherPS = otherPS
	mpmc.thisPlrNum = plrNum
	mpmc.otherPlrNum = otherPlrNum

	plr1.expands = append(plr1.expands, mpmc)
}

type MultiPlayerMissileCollision struct {
	g           *Game
	thisPlr     *Entity
	myPS        *PlayerShooting
	otherPS     *PlayerShooting
	thisPlrNum  int
	otherPlrNum int
}

func (mpmc *MultiPlayerMissileCollision) ExpandUpdate(e *Entity) error {
	for _, otherMissile := range mpmc.otherPS.activePool {
		if CheckCollision(mpmc.thisPlr, otherMissile, 0) {
			mpmc.g.plrHealths[mpmc.thisPlrNum] -= 1
			otherMissile.pos.Y = MissileDeleteYRange
			if mpmc.g.plrHealths[mpmc.thisPlrNum] < 1 {
				mpmc.g.multiWinnerColor = MultiPlayerColors[mpmc.otherPlrNum]
				mpmc.g.ChangeState(StateMultiPlayerOver)
				break
			}
		}
		for _, myMissile := range mpmc.myPS.activePool {
			if CheckCollision(otherMissile, myMissile, 0) {
				myMissile.pos.Y = MissileDeleteYRange
				otherMissile.pos.Y = MissileDeleteYRange
			}
		}
	}
	return nil
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
