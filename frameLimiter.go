package main

import (
	"fmt"
	"time"

	"github.com/faiface/pixel/pixelgl"
)

func NewFrameLimiter() *FrameLimiter {
	return &FrameLimiter{
		millisPerFrame: MillersecondsPerFrame,
		last:           time.Now(),
	}
}

type FrameLimiter struct {
	millisPerFrame int64
	last           time.Time
	frames         uint
	second         <-chan time.Time
}

func (f *FrameLimiter) ShouldDoNextFrame() bool {
	dt := time.Since(f.last).Milliseconds()
	if dt > f.millisPerFrame {
		f.last = time.Now()
		return true
	}
	return false
}

func (f *FrameLimiter) InitFrameCounter() {
	f.frames = 0
	f.second = time.Tick(time.Second)
}

func (f *FrameLimiter) SetTitleWithFPS(win *pixelgl.Window, wcfg *pixelgl.WindowConfig) {
	f.frames++
	select {
	case <-f.second:
		win.SetTitle(fmt.Sprintf("%s | FPS: %d", wcfg.Title, f.frames))
		f.frames = 0
	default:
	}
}
