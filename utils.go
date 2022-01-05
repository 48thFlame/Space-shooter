package main

import (
	"fmt"
	"image"
	"image/color"
	"io/ioutil"
	"math"
	"math/rand"
	"os"
	"time"

	pix "github.com/faiface/pixel"
	pixgl "github.com/faiface/pixel/pixelgl"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"

	"github.com/faiface/pixel/text"
)

const (
	FPS                   = 50 // ! not exact !!
	MillersecondsPerFrame = 1000 / FPS

	WindowWidth  = 1120
	WindowHeight = 864

	StateMenu        = 1
	StatePlaySingal  = 2
	StateSignalOver  = 3
	StateMultiPlayer = 4
	StateQuit        = 5
)

func Initialize(windowTitle string) *Game {
	rand.Seed(time.Now().UnixNano())
	ico, err := LoadPicture("menu/icon.png")
	if err != nil {
		panic(err)
	}
	wcfg := pixgl.WindowConfig{
		Title:  windowTitle,
		Bounds: pix.R(0, 0, WindowWidth, WindowHeight),
		Icon: []pix.Picture{ico},
	}
	win, err := pixgl.NewWindow(wcfg)
	if err != nil {
		panic(err)
	}
	win.SetSmooth(true)
	

	framer := NewFrameCounter()

	g := &Game{wcfg: &wcfg, win: win, framer: framer}
	g.quit = false
	g.state = StateMenu
	g.score = 0
	g.level = 1
	g.tw = NewTextWriter(g)

	return g
}

func LoadPicture(path string) (pix.Picture, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return pix.PictureDataFromImage(img), nil
}

func ChangeXY(a float64) (float64, float64) {
	return -1 * math.Sin(a), math.Cos(a)
}

func Contains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func RemoveFromESlice(s []*Entity, i int) []*Entity {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func EntityR(e *Entity) pix.Rect {
	ewh := e.dim.width / 2 // entity width half
	ehh := e.dim.height / 2
	ex := e.pos.X
	ey := e.pos.Y

	return pix.R(ex-ewh, ey-ehh, ex+ewh, ey+ehh)
}

func CheckCollision(e1, e2 *Entity) bool {
	e1r := EntityR(e1)
	e2r := EntityR(e2)

	return e1r.Intersects(e2r)
}

func EToR(e *Entity) pix.Rect {
	return pix.R(e.pos.X-e.dim.width/2, e.pos.Y-e.dim.height/2, e.pos.X+e.dim.width/2, e.pos.Y+e.dim.height/2)
}

func KeyPressed(e *Entity, key pixgl.Button) bool {
	return e.win.Pressed(key)
}

// // spacer

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

// // spacer

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

// // spacer

func LoadTTF(path string, size float64) (font.Face, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	font, err := truetype.Parse(bytes)
	if err != nil {
		return nil, err
	}

	return truetype.NewFace(font, &truetype.Options{
		Size:              size,
		GlyphCacheEntries: 1,
	}), nil
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

func (tw *TextWriter) WriteText(txt string, pos pix.Vec, color color.Color) {
	tw.txt.Clear()
	tw.txt.Color = color
	tw.txt.WriteString(txt)
	txtMoved := pix.IM.Moved( /*tw.g.win.Bounds().Center().Sub(tw.txt.Bounds().Center())*/ pos)
	tw.txt.Draw(tw.g.win, txtMoved)
}
