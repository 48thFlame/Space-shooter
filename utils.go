package main

import (
	"fmt"
	"image"
	"math"
	"os"
	"time"

	pix "github.com/faiface/pixel"
	pixgl "github.com/faiface/pixel/pixelgl"
)

const (
	FPS                   = 50 // ! not exact !!
	MillersecondsPerFrame = 1000 / FPS

	WindowWidth  = 1120
	WindowHeight = 864
)

func Initialize(windowTitle string) *Game { // (*pgl.WindowConfig, *pgl.Window, *FrameLimiter)
	wcfg := pixgl.WindowConfig{
		Title:  windowTitle,
		Bounds: pix.R(0, 0, WindowWidth, WindowHeight),
	}
	win, err := pixgl.NewWindow(wcfg)
	if err != nil {
		panic(err)
	}
	win.SetSmooth(true)

	framer := NewFrameLimiter()

	g := Game{wcfg: &wcfg, win: win, framer: framer}

	return &g
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
	// return math.Cos(a), -1 * math.Sin(a)
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

// // spacer

func NewFrameLimiter() *FrameLimiter {
	return &FrameLimiter{
		frames: 0,
		second: time.Tick(time.Second),
	}
		// millisPerFrame: MillersecondsPerFrame,
		// last:           time.Now(),
}

type FrameLimiter struct {
	// millisPerFrame int64
	// last           time.Time
	frames         uint
	second         <-chan time.Time
}

// func (f *FrameLimiter) ShouldDoNextFrame() bool {
// 	dt := time.Since(f.last).Milliseconds()
// 	if dt > f.millisPerFrame {
// 		f.last = time.Now()
// 		return true
// 	}
// 	return false
// }

// func (f *FrameLimiter) InitFrameCounter() {
// 	f.frames = 0
// 	f.second = time.Tick(time.Second)
// }

func (f *FrameLimiter) SetTitleWithFPS(win *pixgl.Window, wcfg *pixgl.WindowConfig) {
	f.frames++
	select {
	case <-f.second:
		win.SetTitle(fmt.Sprintf("%s | FPS: %d", wcfg.Title, f.frames))
		f.frames = 0
	default:
	}
}
