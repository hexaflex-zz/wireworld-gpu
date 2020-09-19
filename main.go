package main

import (
	"runtime"

	"github.com/go-gl/glfw/v3.3/glfw"
)

// Make sure main() and all openGL related stuff runs in the main thread.
func init() {
	runtime.LockOSThread()
}

func main() {
	var app Application

	app.Initialize()
	defer app.Release()

	for !app.window.ShouldClose() {
		app.Update()
		app.Draw()
		glfw.PollEvents()
	}
}
