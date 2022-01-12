package main

import (
	"fmt"
	"image/color"
	"time"

	pix "github.com/faiface/pixel"
	pixgl "github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
)

type Game struct {
	quit             bool
	state            int
	background       *Entity
	initStateMap     map[int]func(g *Game)
	wcfg             *pixgl.WindowConfig
	win              *pixgl.Window
	framer           *FrameCounter
	entities         []*Entity
	tw               *TextWriter
	score            int
	level            int
	plrHealths       map[int]int
	multiWinnerColor string
}

func (g *Game) ChangeState(state int) {
	g.initStateMap[state](g)
	g.state = state
}

func (g *Game) WinClear() {
	g.win.Clear(color.RGBA{30, 30, 30, 255})
	g.background.Draw()
}

func (g *Game) WinUpdate() {
	g.win.Update()
	time.Sleep(time.Millisecond * MillersecondsPerFrame)
}

func (g *Game) AddEntity(e *Entity) {
	g.entities = append(g.entities, e)
}

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

	e.rot = 0
	e.pos = pix.Vec{X: 0, Y: 0}
	e.dim = Dimension{width: width, height: height}
	e.img = img
	e.spr = spr
	e.win = win

	return &e
}

type Entity struct {
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
	mat := pix.IM
	mat = mat.Moved(e.pos)
	mat = mat.Rotated(e.pos, e.rot)

	e.spr.Draw(e.win, mat)
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

func NewFrameCounter() *FrameCounter {
	return &FrameCounter{
		frames: 0,
		second: time.Tick(time.Second),
	}
}

type FrameCounter struct {
	frames uint
	second <-chan time.Time
}

func (f *FrameCounter) SetTitleWithFPS(win *pixgl.Window, wcfg *pixgl.WindowConfig) {
	f.frames++
	select {
	case <-f.second:
		win.SetTitle(fmt.Sprintf("%s | FPS: %d", wcfg.Title, f.frames))
		f.frames = 0
	default:
	}
}

func NewTextWriter(g *Game) *TextWriter {
	tw := &TextWriter{}

	tw.g = g
	font, err := LoadTTF("arial.ttf", 48)
	if err != nil {
		panic(err)
	}

	tw.atlas = text.NewAtlas(font, text.ASCII)
	tw.txt = text.New(pix.V(0, 0), tw.atlas)

	return tw
}

type TextWriter struct {
	g     *Game
	atlas *text.Atlas
	txt   *text.Text
}

func (tw *TextWriter) WriteText(txt string, pos pix.Vec, color color.Color, size float64) {
	tw.txt.Clear()
	tw.txt.Color = color
	tw.txt.WriteString(txt)
	txtMoved := pix.IM.Moved(pos).Scaled(pos, size)
	tw.txt.Draw(tw.g.win, txtMoved)
}
