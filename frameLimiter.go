package main

import "time"

func NewFrameLimiter() *FrameLimiter {
	return &FrameLimiter{
		millisPerFrame: MillersecondsPerFrame,
		last:           time.Now(),
	}
}

type FrameLimiter struct {
	millisPerFrame int64
	last           time.Time
}

func (f *FrameLimiter) ShouldDoNextFrame() bool {
	dt := time.Since(f.last).Milliseconds()
	if dt > f.millisPerFrame {
		f.last = time.Now()
		return true
	}
	return false
}
