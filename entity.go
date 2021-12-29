package main

import (
	"fmt"

	pix "github.com/faiface/pixel"
	pixgl "github.com/faiface/pixel/pixelgl"
)

func NewEntity(pictureFile string, win *pixgl.Window) *Entity {
	e := Entity{}

	img, err := LoadPicture(pictureFile)
	if err != nil {
		panic(fmt.Errorf("can't load picture %v, to create a new Entity: %v", pictureFile, err))
	}
	spr := pix.NewSprite(img, img.Bounds())
	rect := img.Bounds()
	width := rect.Max.X - rect.Min.X
	height := rect.Max.Y - rect.Min.Y

	e.active = true
	e.rot = 0
	e.pos = pix.Vec{X: 0, Y: 0}
	e.dim = Dimension{width: width, height: height}
	e.img = img
	e.spr = spr
	e.win = win

	return &e
}

type Entity struct {
	active  bool
	rot     float64       // rotation
	pos     pix.Vec       // postition
	dim     Dimension     // dimensions
	img     pix.Picture   // image
	spr     *pix.Sprite   // sprite
	win     *pixgl.Window // window
	expands []Expantion   // expansions
}

type Dimension struct {
	width, height float64
}

type Expantion interface {
	ExpandUpdate(*Entity) error
}

func (e *Entity) Draw() {
	if e.active {
		mat := pix.IM
		mat = mat.Moved(e.pos)
		mat = mat.Rotated(e.pos, e.rot)

		e.spr.Draw(e.win, mat)
	}
}

func (e *Entity) Update() error {
	for _, expansion := range e.expands {
		err := expansion.ExpandUpdate(e)
		if err != nil {
			return err
		}
	}
	return nil
}

func (e *Entity) KeyPressed(key pixgl.Button) bool {
	return e.win.Pressed(key)
}
