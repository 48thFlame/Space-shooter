package main

import (
	"image"
	"os"
	"time"

	pix "github.com/faiface/pixel"
	pgl "github.com/faiface/pixel/pixelgl"
)

const (
	FPS = 50  // ! not exact !!
	MillersecondsPerFrame = 1000 / FPS

	WindowWidth = 624
	WindowHeight = 912

	PlayerSpeed = 7
	PlayerCooldown = time.Millisecond * 125
	PlayerMissileSpeed = 12
	PlayerMissilePoolSize = 80
)

func Initialize(windowTitle string) (*pgl.WindowConfig, *pgl.Window, *FrameLimiter) {
	wcfg := pgl.WindowConfig{
		Title:  windowTitle,
		Bounds: pix.R(0, 0, WindowWidth, WindowHeight),
	}
	win, err := pgl.NewWindow(wcfg)
	if err != nil {
		panic(err)
	}

	framer := NewFrameLimiter()

	return &wcfg, win, framer
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
