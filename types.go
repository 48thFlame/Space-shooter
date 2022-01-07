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
	quit     bool
	state    int
	wcfg     *pixgl.WindowConfig
	win      *pixgl.Window
	framer   *FrameCounter
	entities []*Entity
	tw       *TextWriter
	score    int
	level    int
}

func (g *Game) WinClear() {
	g.win.Clear(color.RGBA{30, 30, 30, 255})
}

func (g *Game) WinUpdate() {
	g.win.Update()
	time.Sleep(time.Millisecond * MillersecondsPerFrame)
}

func (g *Game) InitSignalPlayerGame() {
	g.entities = []*Entity{}
	g.score = 0
	g.level = 1
	plr := NewSingalPlayer(g)

	g.AddEntity(plr)
}

func (g *Game) SignalPlayerGame() {
	g.WinClear()

	g.tw.WriteText(fmt.Sprint(g.score), pix.V(WindowWidth-128, WindowHeight-82), color.RGBA{220, 220, 220, 255})

	for _, e := range g.entities {
		e.Update()
		e.Draw()
	}

	g.framer.SetTitleWithFPS(g.win, g.wcfg)

	g.WinUpdate()
}

func (g *Game) InitEndSingalPlayer() {
	g.entities = []*Entity{}
	menu := NewEntity("menu/game over.png", g.win)
	menu.pos = pix.V(WindowWidth/2, WindowHeight/2)

	button := NewButton("menu/button 4.png", pix.V(WindowWidth/2, WindowHeight/2-300), g, StateMenu)

	g.AddEntity(menu)

	g.AddEntity(button)
}

func (g *Game) EndSingalPlayer() {
	g.WinClear()
	for _, e := range g.entities {
		e.Update()
		e.Draw()
	}

	g.tw.WriteText(
		fmt.Sprint(g.score),
		pix.V(750, 453),
		color.RGBA{255, 255, 255, 255},
	)

	g.WinUpdate()
}

func (g *Game) InitMenu() {
	bmh := 2.75

	g.entities = []*Entity{}
	background := NewEntity("menu/menu.png", g.win)
	background.pos = pix.V(WindowWidth/2, WindowHeight/2)

	mouse := NewEntity("menu/mouse.png", g.win)
	mouse.expands = append(mouse.expands, &MouseFollower{})

	b1 := NewButton("menu/button 1.png", pix.V(WindowWidth/2, WindowHeight/bmh), g, StatePlaySingal)
	b2 := NewButton("menu/button 2.png", pix.V(WindowWidth/2, WindowHeight/bmh-56), g, StateMultiPlayer)
	b3 := NewButton("menu/button 3.png", pix.V(WindowWidth/2, WindowHeight/bmh-112), g, StateQuit)

	g.AddEntity(background)

	g.AddEntity(b1)
	g.AddEntity(b2)
	g.AddEntity(b3)

	g.AddEntity(mouse)
}

func (g *Game) Menu() {
	g.WinClear()
	for _, e := range g.entities {
		e.Update()
		e.Draw()
	}
	g.WinUpdate()
}

func (g *Game) AddEntity(e *Entity) {
	g.entities = append(g.entities, e)
}

type MouseFollower struct{}

func (mf *MouseFollower) ExpandUpdate(e *Entity) error {
	e.pos = e.win.MousePosition()
	return nil
}

func NewButton(filePath string, pos pix.Vec, g *Game, id int) *Entity {
	e := NewEntity(filePath, g.win)

	e.pos = pos
	mb := &UIButton{}
	mb.id = id
	mb.g = g

	e.expands = append(e.expands, mb)

	return e
}

type UIButton struct {
	id int
	g  *Game
	// label string
}

func (mb *UIButton) ExpandUpdate(e *Entity) error {
	if mb.g.win.Pressed(pixgl.MouseButtonLeft) {
		mousePos := e.win.MousePosition()
		if EToR(e, 0).Contains(mousePos) {
			if mb.id == StatePlaySingal {
				mb.g.InitSignalPlayerGame()
			} else if mb.id == StateMenu {
				mb.g.InitMenu()
			}
			mb.g.state = mb.id
		}
	}
	return nil
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

func (tw *TextWriter) WriteText(txt string, pos pix.Vec, color color.Color) {
	tw.txt.Clear()
	tw.txt.Color = color
	tw.txt.WriteString(txt)
	txtMoved := pix.IM.Moved(pos)
	tw.txt.Draw(tw.g.win, txtMoved)
}
