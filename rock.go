package main

import (
	"fmt"
	"math"
	"math/rand"
	"sort"

	pix "github.com/faiface/pixel"
)

const (
	MissileDeleteYRange                 = -10000 // the y that if a missile are below should delete them
	NumOfRockPerLevel                   = 4
	RockPlayerCollsionForgivenessAmount = 8
	RockMaxSpeed                        = 4.5
	RockMinSpeed                        = 1.5
	RockDestroyScoreNum                 = 10
	RockOffEdgeCorrection               = 56
)

func NewRockController(plr *Entity, g *Game) *PlayerRockController {
	prc := &PlayerRockController{}

	prc.g = g
	prc.GenLevel(g, plr)

	return prc
}

// expansion to go on player
type PlayerRockController struct {
	pool []*Entity
	g    *Game
}

func (prc *PlayerRockController) ExpandUpdate(plr *Entity) error {
	for _, rock := range prc.pool {
		rock.Update()
		rock.Draw()
	}

	return nil
}

func NewRock(prc *PlayerRockController, plr *Entity, size int, g *Game) *Entity {
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
	re.g = prc.g

	rock.expands = append(rock.expands, &re)

	return rock
}

// expansion to go on the rock entity
type RockExpand struct {
	plr   *Entity
	speed float64
	size  int
	g     *Game
}

func (re *RockExpand) Moverock(rock *Entity) {
	xv, yv := ChangeXY(rock.rot)

	xv *= re.speed
	yv *= re.speed

	rock.pos = rock.pos.Add(pix.V(xv, yv))

	// if off edge should go to other side
	if rock.pos.X+RockOffEdgeCorrection > WindowWidth+rock.dim.width {
		rock.pos.X = 0 + RockOffEdgeCorrection
	} else if rock.pos.X-RockOffEdgeCorrection < 0-rock.dim.width {
		rock.pos.X = WindowWidth - RockOffEdgeCorrection
	}
	if rock.pos.Y+RockOffEdgeCorrection > WindowHeight+rock.dim.height {
		rock.pos.Y = 0 + RockOffEdgeCorrection
	} else if rock.pos.Y-RockOffEdgeCorrection < 0-rock.dim.height {
		rock.pos.Y = WindowHeight - RockOffEdgeCorrection
	}
}

func (re *RockExpand) ExpandUpdate(rock *Entity) error {
	re.Moverock(rock)

	if CheckCollision(re.plr, rock, RockPlayerCollsionForgivenessAmount) {
		re.g.ChangeState(StateSingalPlayerOver)
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
			if CheckCollision(missile, rock, RockPlayerCollsionForgivenessAmount) {
				pmrc.g.score += RockDestroyScoreNum
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
		prc.pool = append(prc.pool, NewRock(prc, plr, rand.Intn(3)+1, prc.g))
	}
}
