package main

import (
	"image"
	"io/ioutil"
	"math"
	"math/rand"
	"os"
	"time"

	pix "github.com/faiface/pixel"
	pixgl "github.com/faiface/pixel/pixelgl"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

const (
	FPS                   = 50 // ! not exact !!
	MillersecondsPerFrame = 1000 / FPS

	WindowWidth  = 1120
	WindowHeight = 864

	BackgroundFilePath = "menu/background.png"

	StateQuit             = 0
	StateMenu             = 1
	StateSingalPlayer     = 2
	StateSingalPlayerOver = 3
	StateHighScore        = 4
	StateMultiPlayer      = 5
	StateMultiPlayerOver  = 6
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
		Icon:   []pix.Picture{ico},
	}
	win, err := pixgl.NewWindow(wcfg)
	if err != nil {
		panic(err)
	}
	win.SetSmooth(true)

	framer := NewFrameCounter()

	background := NewEntity(BackgroundFilePath, win)
	background.pos = pix.V(WindowWidth/2, WindowHeight/2)

	g := &Game{wcfg: &wcfg, win: win, framer: framer}
	g.quit = false
	g.background = background
	g.score = 0
	g.level = 1
	g.plrHealths = MultiHealths
	g.tw = NewTextWriter(g)
	g.initStateMap = map[int]func(*Game){
		StateQuit:             InitQuitState,
		StateMenu:             InitMenu,
		StateSingalPlayer:     InitSignalPlayerState,
		StateSingalPlayerOver: InitSingalPlayerOverState,
		// StateHighScore:        InitHighScoreState,
		StateMultiPlayer:     InitMultiPlayerState,
		StateMultiPlayerOver: InitMultiPlayerOverState,
	}
	g.ChangeState(StateMenu)

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

func CheckCollision(e1, e2 *Entity, foreignnessAmount float64) bool {
	e1r := EToR(e1, foreignnessAmount)
	e2r := EToR(e2, foreignnessAmount)

	return e1r.Intersects(e2r)
}

func EToR(e *Entity, foreignnessAmount float64) pix.Rect {
	return pix.R(
		e.pos.X-e.dim.width/2+foreignnessAmount,
		e.pos.Y-e.dim.height/2+foreignnessAmount,
		e.pos.X+e.dim.width/2-foreignnessAmount,
		e.pos.Y+e.dim.height/2-foreignnessAmount,
	)
}

func KeyPressed(e *Entity, key pixgl.Button) bool {
	return e.win.Pressed(key)
}
