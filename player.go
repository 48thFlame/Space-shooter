package main

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
	"time"

	pix "github.com/faiface/pixel"
	pixgl "github.com/faiface/pixel/pixelgl"
)

const (
	PlayerAccSpeed                 = 0.66
	PlayerMaxSpeed                 = 16
	PlayerPercentOfSpeedShouldSlow = 0.02
	PlayerRotSpeed                 = 0.17
	PlayerCooldown                 = time.Millisecond * 125
	PlayerMissileSpeed             = 25
	PlayerMissilePoolSize          = 50
	NumOfRockPerLevel              = 4
	RockMaxSpeed                   = 4.5
	RockMinSpeed                   = 1.5
	MissileDeleteYRange            = -10000 // the y that if a missile are below should delete them
	RockDestroyScoreNum            = 10
)

func NewPlayer(g *Game) *Entity {
	plr := NewEntity("sprites/player.png", g.win)

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

func NewRockController(plr *Entity, g *Game) *PlayerRockController {
	prc := &PlayerRockController{}

	prc.GenLevel(g, plr)

	return prc
}

// expansion to go on player
type PlayerRockController struct {
	pool []*Entity
}

func (prc *PlayerRockController) ExpandUpdate(plr *Entity) error {
	for _, rock := range prc.pool {
		rock.Update()
		rock.Draw()
	}

	return nil
}

func NewRock(prc *PlayerRockController, plr *Entity, size int) *Entity {
	rock := NewEntity(fmt.Sprintf("sprites/rock%v.png", size), plr.win)
	rock.rot = rand.Float64() * (math.Pi*4 - 0)

	rockXOptions := []float64{0, float64(WindowWidth)}
	rockYOptions := []float64{0, float64(WindowHeight)}

	rock.pos.X = rockXOptions[rand.Intn(2)]
	rock.pos.Y = rockYOptions[rand.Intn(2)]

	re := RockExpand{}
	re.speed = RockMinSpeed + rand.Float64()*(RockMaxSpeed-RockMinSpeed)
	re.size = size
	re.plr = plr

	rock.expands = append(rock.expands, &re)

	return rock
}

// expansion to go on the rock entity
type RockExpand struct {
	plr   *Entity
	speed float64
	size  int
}

func (re *RockExpand) Moverock(rock *Entity) {
	xv, yv := ChangeXY(rock.rot)

	xv *= re.speed
	yv *= re.speed

	rock.pos = rock.pos.Add(pix.V(xv, yv))

	// if off edge should go to other side
	if rock.pos.X > WindowWidth+rock.dim.width {
		rock.pos.X = 0
	} else if rock.pos.X < 0-rock.dim.width {
		rock.pos.X = WindowWidth
	}
	if rock.pos.Y > WindowHeight+rock.dim.height {
		rock.pos.Y = 0
	} else if rock.pos.Y < 0-rock.dim.height {
		rock.pos.Y = WindowHeight
	}
}

func (re *RockExpand) ExpandUpdate(rock *Entity) error {
	re.Moverock(rock)

	if CheckCollision(re.plr, rock) {
		fmt.Println("Game over!")
	}

	return nil
}

type PlayerMissileRockCollision struct {
	ps  *PlayerShooting
	prc *PlayerRockController
	g   *Game
	plr *Entity
}

func (pmrc *PlayerMissileRockCollision) ExpandUpdate(rock *Entity) error {
	var rockIRemove []int
	var rocksRemoved int

	for _, missile := range pmrc.ps.activePool {
		for rockI, rock := range pmrc.prc.pool {
			if CheckCollision(missile, rock) {
				pmrc.g.score += RockDestroyScoreNum
				fmt.Println(pmrc.g.score)
				rockIRemove = append(rockIRemove, rockI)
				missile.pos.Y = MissileDeleteYRange
			}
		}
	}

	// remove all the hit rocks
	sort.Ints(rockIRemove)
	for _, i := range rockIRemove {
		pmrc.prc.pool = RemoveFromESlice(pmrc.prc.pool, i-rocksRemoved)
		rocksRemoved++
	}

	if len(pmrc.prc.pool) == 0 {
		pmrc.g.level += 1
		pmrc.prc.GenLevel(pmrc.g, pmrc.plr)
		return nil
	}

	return nil
}

func (prc *PlayerRockController) GenLevel(g *Game, plr *Entity) {
	for i := 0; i < NumOfRockPerLevel*g.level; i++ {
		prc.pool = append(prc.pool, NewRock(prc, plr, rand.Intn(2)+1))
	}
}
