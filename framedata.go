package main

import "time"

// FrameData keeps track of frame timing data.
type FrameData struct {
	lastUpdate  time.Time
	noLongerNow time.Time
	framecount  int
	framerate   int
	frameDelta  time.Duration
}

// FrameDelta returns the time passed since the last update.
func (f *FrameData) FrameDelta() time.Duration {
	return f.frameDelta
}

// Framerate returns the current framerate value.
func (f *FrameData) Framerate() int {
	return f.framerate
}

// Update updates FrameData state.
func (f *FrameData) Update(now time.Time) {
	f.framecount++
	f.frameDelta = now.Sub(f.noLongerNow)
	f.noLongerNow = now

	if now.Sub(f.lastUpdate) >= time.Second {
		f.framerate = f.framecount
		f.lastUpdate = now
		f.framecount = 0
	}
}
