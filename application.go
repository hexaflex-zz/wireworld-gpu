package main

import (
	"fmt"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/go-gl/gl/v4.2-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/hexaflex/wireworld-gpu/math"
	"github.com/pkg/errors"
)

// Application defines application state.
type Application struct {
	config         *Config
	window         *glfw.Window
	simulation     *Simulation
	display        *SimulationDisplay
	mouse          math.Vec2
	mouseDelta     math.Vec2
	titleUpdated   time.Time
	stepTime       time.Time
	stepInterval   time.Duration
	stepMultiplier int
	clockCycles    uint64
	uboShared      uint32
	running        bool
	glInitialized  bool
}

// Initialize initializes the window and openGL.
func (a *Application) Initialize() {
	var err error

	a.config = parseArgs()
	a.stepInterval = time.Millisecond * 10
	a.stepMultiplier = 1

	log.Println(Version())
	a.check(glfw.Init())
	a.check(a.setWindowMode(a.config.Width, a.config.Height, a.config.Fullscreen))
	a.check(gl.Init())
	a.glInitialized = true

	glver := gl.GoStr(gl.GetString(gl.VERSION))
	log.Println("using OpenGL version:", glver)

	gl.Disable(gl.DEPTH_TEST)
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	// layout(std140, binding = 0) uniform Shared {
	// 	 mat4 View;
	// 	 mat4 Projection;
	// };
	const sizeofMat4 = len(math.Mat4{}) * 4
	const structSize = 2 * sizeofMat4
	gl.GenBuffers(1, &a.uboShared)
	gl.BindBuffer(gl.UNIFORM_BUFFER, a.uboShared)
	gl.BufferData(gl.UNIFORM_BUFFER, structSize, nil, gl.DYNAMIC_DRAW)
	gl.BindBufferRange(gl.UNIFORM_BUFFER, 0, a.uboShared, 0, structSize)
	gl.BindBuffer(gl.UNIFORM_BUFFER, 0)

	a.simulation, err = LoadSimulation(a.config.Input, &a.config.Palette)
	a.check(err)

	displayShader, err := DisplayShader.Compile()
	a.check(err)

	w, h := a.window.GetFramebufferSize()

	a.display = NewSimulationDisplay(displayShader)
	a.display.SetSize(a.simulation.Size())
	a.display.SetPalette(&a.config.Palette)
	a.display.Center(math.Vec2{float32(w), float32(h)})

	// Force resize call now that components have been initialized.
	a.framebufferSizeCallback(a.window, w, h)
}

// Release clears resources.
func (a *Application) Release() {
	gl.DeleteBuffers(1, &a.uboShared)

	if a.simulation != nil {
		a.simulation.Release()
		a.simulation = nil
	}

	if a.display != nil {
		a.display.Release()
		a.display = nil
	}

	if a.window != nil {
		a.window.SetKeyCallback(nil)
		a.window.SetFramebufferSizeCallback(nil)
		a.window.SetScrollCallback(nil)
		a.window.SetCursorPosCallback(nil)
		a.window.Destroy()
		a.window = nil
	}

	glfw.Terminate()
}

// Update updates application state and handles input.
func (a *Application) Update() {
	now := time.Now()

	if now.Sub(a.titleUpdated) >= time.Second {
		a.titleUpdated = now
		state := "stopped"
		if a.running {
			state = "running"
		}
		text := fmt.Sprintf(
			"%s - [%s] clock: %s",
			Version(),
			state,
			a.clockFrequency(),
		)
		a.window.SetTitle(text)
	}

	if a.running && now.Sub(a.stepTime) >= a.stepInterval {
		a.stepTime = now
		a.clockCycles += uint64(a.stepMultiplier)
		a.simulation.Step(a.stepMultiplier)
	}
}

// Draw renders the scene.
func (a *Application) Draw() {
	w, h := a.window.GetFramebufferSize()
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.Viewport(0, 0, int32(w), int32(h))
	a.display.Draw(a.simulation)
	a.window.SwapBuffers()
}

func (a *Application) framebufferSizeCallback(window *glfw.Window, width, height int) {
	if !a.glInitialized {
		return
	}

	gl.Viewport(0, 0, int32(width), int32(height))
	a.updateUniformBlock()
}

func (a *Application) cursorPosCallback(window *glfw.Window, x, y float64) {
	pos := math.Vec2{float32(x), float32(y)}
	a.mouseDelta = a.mouse.Sub(pos)
	a.mouse = pos

	if a.display != nil && a.window.GetKey(glfw.KeySpace) != glfw.Release {
		a.display.Scroll(a.mouseDelta)
	}
}

func (a *Application) scrollCallback(window *glfw.Window, x, y float64) {
	if a.display != nil {
		a.display.Zoom(float32(y), a.mouse)
	}
}

func (a *Application) keyCallback(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if action != glfw.Press {
		return
	}

	switch key {
	case glfw.KeyEscape:
		window.SetShouldClose(true)
	case glfw.KeyF1:
		a.saveState()
	case glfw.KeyF2:
		a.loadState()
	case glfw.KeyF5:
		a.reload()
	case glfw.KeyF11:
		w, h := a.window.GetFramebufferSize()
		a.setWindowMode(w, h, !a.config.Fullscreen)
	case glfw.KeyC:
		w, h := a.window.GetFramebufferSize()
		a.display.Center(math.Vec2{float32(w), float32(h)})
	case glfw.KeyQ:
		a.running = !a.running
	case glfw.KeyE:
		a.simulation.Step(1)
	case glfw.KeyW:
		a.increaseClockspeed()
	case glfw.KeyS:
		a.decreaseClockspeed()
	}
}

// decreaseClockspeed slows the clock down.
func (a *Application) decreaseClockspeed() {
	// Using just a time interval, we can't go above 1kHz.
	//
	// To break this limit, we use an additional step multiplier
	// which is always 1, until we hit the 1kHz limit.
	//
	// The multiplier simply increases the number of times the
	// step function is run during a single cycle. With this, we
	// can increase execution speed to up to 80-90 Khz or more,
	// depending on the GPU being used.

	if a.stepMultiplier > 1 {
		a.stepMultiplier /= 10
	} else {
		a.stepInterval *= 10
		if a.stepInterval > time.Second {
			a.stepInterval = time.Nanosecond
		}
	}
}

// increaseClockspeed speeds up the execution.
func (a *Application) increaseClockspeed() {
	// Using just a time interval, we can't go above 1kHz.
	//
	// To break this limit, we use an additional step multiplier
	// which is always 1, until we hit the 1kHz limit.
	//
	// The multiplier simply increases the number of times the
	// step function is run during a single cycle. With this, we
	// can increase execution speed to up to 80-90 Khz or more,
	// depending on the GPU being used.

	if a.stepInterval <= time.Millisecond {
		a.stepMultiplier *= 10
	} else {
		a.stepInterval /= 10
		if a.stepInterval < time.Millisecond {
			a.stepInterval = time.Millisecond
		}
	}
}

// clockFrequency returns the current clock frequency in a prettyfied string.
// Called once per second.
func (a *Application) clockFrequency() string {
	freq := float32(a.clockCycles)
	a.clockCycles = 0

	switch {
	case freq > 1e9:
		return fmt.Sprintf("%.2f Ghz", freq/1e9)
	case freq > 1e6:
		return fmt.Sprintf("%.2f Mhz", freq/1e6)
	case freq > 1e3:
		return fmt.Sprintf("%.2f Khz", freq/1e3)
	default:
		return fmt.Sprintf("%.0f Hz", freq)
	}
}

// check writes the given error to stderr and exits the program, if the error is not nil.
func (a *Application) check(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		a.Release()
		os.Exit(1)
	}
}

// setWindowMode sets the window mode. If no window exists yet, a new one
// is created. Returns an error if this process fails. If a window already
// exists, this always returns nil.
func (a *Application) setWindowMode(width, height int, fullscreen bool) error {
	var x, y int

	a.config.Fullscreen = fullscreen
	monitor := glfw.GetPrimaryMonitor()
	mode := monitor.GetVideoMode()

	if !fullscreen {
		x = (mode.Width / 2) - (width / 2)
		y = (mode.Height / 2) - (height / 2)
		monitor = nil // must be nil if we want windowed mode.
	}

	if a.window != nil {
		a.window.SetMonitor(monitor, x, y, width, height, mode.RefreshRate)
		return nil
	}

	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 2)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	glfw.WindowHint(glfw.ScaleToMonitor, glfw.True)
	glfw.WindowHint(glfw.Focused, glfw.True)
	glfw.WindowHint(glfw.Decorated, glfw.True)
	glfw.WindowHint(glfw.Maximized, glfw.False)
	glfw.WindowHint(glfw.Resizable, glfw.True)

	if !fullscreen {
		glfw.WindowHint(glfw.Visible, glfw.False)
	}

	var err error
	a.window, err = glfw.CreateWindow(width, height, Version(), monitor, nil)
	if err != nil {
		return errors.Wrapf(err, "glfw.CreateWindow failed")
	}

	a.window.SetKeyCallback(a.keyCallback)
	a.window.SetFramebufferSizeCallback(a.framebufferSizeCallback)
	a.window.SetScrollCallback(a.scrollCallback)
	a.window.SetCursorPosCallback(a.cursorPosCallback)
	a.window.MakeContextCurrent()
	a.window.SetSize(width, height)
	a.window.SetPos(x, y)
	a.window.Show()

	glfw.SwapInterval(0)

	return nil
}

// updateUniformBlock updates matrix information in the shared uniform block.
func (a *Application) updateUniformBlock() {
	w, h := a.window.GetFramebufferSize()
	p := math.Ortho2D(0, float32(w), float32(h), 0)
	v := math.Ident4()

	const sizeofMat4 = len(math.Mat4{}) * 4
	gl.BindBuffer(gl.UNIFORM_BUFFER, a.uboShared)
	gl.BufferSubData(gl.UNIFORM_BUFFER, 0, sizeofMat4, gl.Ptr(v[:]))
	gl.BufferSubData(gl.UNIFORM_BUFFER, sizeofMat4, sizeofMat4, gl.Ptr(p[:]))
	gl.BindBuffer(gl.UNIFORM_BUFFER, 0)
}

// reload reloads the original input image from disk.
func (a *Application) reload() {
	var err error

	log.Println("reloading", a.config.Input)

	a.simulation, err = LoadSimulation(a.config.Input, &a.config.Palette)
	if err != nil {
		log.Println("load failed:", err)
	}
}

// saveState writes the current simulation state as a PNG file.
func (a *Application) saveState() {
	stamp := time.Now().UnixNano()
	dir, name := filepath.Split(a.config.Input)
	name = strings.Replace(name, filepath.Ext(name), "", -1)
	file := filepath.Join(dir, fmt.Sprintf("%d.%s.png", stamp, name))

	log.Println("saving state file", file)

	fd, err := os.Create(file)
	if err != nil {
		log.Println("failed to create state:", err)
		return
	}

	img := a.simulation.Image(&a.config.Palette)
	if err = png.Encode(fd, img); err != nil {
		log.Println("failed to encode image:", err)
		_ = fd.Close()
		return
	}

	if err = fd.Close(); err != nil {
		log.Println("failed to save state:", err)
	}
}

// loadState loads the latest state of the input image from disk.
func (a *Application) loadState() {
	dir, name := filepath.Split(a.config.Input)

	// Find all files matching the input.
	stamps := findStateFiles(dir, name)
	if len(stamps) == 0 {
		log.Println("no state file found")
		a.reload()
		return
	}

	sort.Slice(stamps, func(i, j int) bool {
		return stamps[i] > stamps[j]
	})

	file := fmt.Sprintf("%d.%s", stamps[0], name)
	file = filepath.Join(dir, file)

	log.Println("loading state", file)

	var err error
	a.simulation, err = LoadSimulation(file, &a.config.Palette)
	if err != nil {
		log.Println("failed to load state:", err)
	}
}

// findStateFiles returns all files from the give directory which
// have a timestamp prefix. It returns those timestamps.
func findStateFiles(dir, name string) []int64 {
	fd, err := os.Open(dir)
	if err != nil {
		return nil
	}

	files, err := fd.Readdirnames(-1)
	fd.Close()
	if err != nil {
		return nil
	}

	out := make([]int64, 0, len(files))
	for _, file := range files {
		if !strings.HasSuffix(file, "."+name) {
			continue
		}

		index := strings.Index(file, ".")
		left := file[:index]
		stamp, err := strconv.ParseInt(left, 10, 64)
		if err != nil {
			continue
		}
		out = append(out, stamp)
	}

	return out
}
