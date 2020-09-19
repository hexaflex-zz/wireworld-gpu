package main

import (
	"image/color"

	"github.com/go-gl/gl/v4.2-core/gl"
	math "github.com/hexaflex/glmath"
)

// Zoom limits for the simulation display.
const (
	MinZoom     = 1
	MaxZoom     = 100
	DefaultZoom = 5
)

// SimulationDisplay is a textured quad that renders the current state of
// a simulation using a given color palette.
type SimulationDisplay struct {
	transform      *math.Transform
	shader         Shader
	zoomFactor     float32
	vao            uint32
	vbo            uint32
	transformDirty bool
}

// NewSimulationDisplay creates a new, blank Display.
func NewSimulationDisplay(shader Shader) *SimulationDisplay {
	var d SimulationDisplay
	var verts = []float32{
		// x,y,u,v
		-0.5, -0.5, 0, 0,
		0.5, -0.5, 1, 0,
		-0.5, 0.5, 0, 1,
		0.5, -0.5, 1, 0,
		0.5, 0.5, 1, 1,
		-0.5, 0.5, 0, 1}

	d.transformDirty = true
	d.transform = math.NewTransform()
	d.shader = shader
	d.SetZoom(DefaultZoom)

	gl.GenVertexArrays(1, &d.vao)
	gl.BindVertexArray(d.vao)

	gl.GenBuffers(1, &d.vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, d.vbo)
	gl.EnableVertexAttribArray(0)
	gl.EnableVertexAttribArray(1)
	gl.VertexAttribPointer(0, 2, gl.FLOAT, false, 4*4, gl.PtrOffset(0))
	gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 4*4, gl.PtrOffset(2*4))
	gl.BufferData(gl.ARRAY_BUFFER, len(verts)*4, gl.Ptr(verts), gl.STATIC_DRAW)

	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindVertexArray(0)
	return &d
}

// Release cleans up resources.
func (d *SimulationDisplay) Release() {
	gl.DeleteBuffers(1, &d.vbo)
	gl.DeleteVertexArrays(1, &d.vao)
}

// Zoom zooms in/out of the map by the given relative offset.
// Focus is used as the point on which the zoom focuses.
func (d *SimulationDisplay) Zoom(delta float32, focus math.Vec2) {
	if math.FloatEqual(delta, 0) {
		return
	}

	// Zooming needs to center on the given focal point.
	// For this to work, we figure out which position on
	// the Display is currently under the cursor. Then we
	// perform the zoom and scroll back to that position.

	origin := d.transform.Translate
	oldZoom := d.zoomFactor

	focus = focus.Sub(origin)
	d.zoomFactor = math.Clamp(oldZoom+delta, MinZoom, MaxZoom)

	xy1 := focus.DivScalar(oldZoom)
	xy2 := focus.Sub(xy1.MulScalar(d.zoomFactor).Sub(origin))

	d.SetScroll(xy2)
}

// SetZoom sets the map's zoom level.
func (d *SimulationDisplay) SetZoom(z float32) {
	d.zoomFactor = math.Clamp(z, MinZoom, MaxZoom)
}

// Center centers the simulation in the given viewport.
func (d *SimulationDisplay) Center(view math.Vec2) {
	size := d.transform.Scale.MulScalar(0.5)
	pos := view.MulScalar(0.5)
	d.SetScroll(pos.Sub(size))
}

// Scroll moves the scroll origin by the given relative offset.
func (d *SimulationDisplay) Scroll(offset math.Vec2) {
	d.transform.Translate = d.transform.Translate.Sub(offset)
	d.transformDirty = true
}

// SetScroll sets the scroll origin to the given position.
func (d *SimulationDisplay) SetScroll(pos math.Vec2) {
	d.transform.Translate = pos
	d.transformDirty = true
}

// SetSize sets the size of the display.
func (d *SimulationDisplay) SetSize(size math.Vec2) {
	d.transform.Scale = size
	d.transformDirty = true
}

// SetPalette sets the color palette used to render the simulation.
func (d *SimulationDisplay) SetPalette(pal *Palette) {
	toVec4 := func(c color.RGBA) math.Vec4 {
		return math.Vec4{
			float32(c.R) / 255,
			float32(c.G) / 255,
			float32(c.B) / 255,
			float32(c.A) / 255,
		}
	}

	d.shader.Use()
	d.shader.SetUniformVec4("PalEmpty", toVec4(pal.Empty))
	d.shader.SetUniformVec4("PalWire", toVec4(pal.Wire))
	d.shader.SetUniformVec4("PalHead", toVec4(pal.Head))
	d.shader.SetUniformVec4("PalTail", toVec4(pal.Tail))
	d.shader.Unuse()
}

// Bindable defines an object with a bindable texture.
type Bindable interface {
	Bind()
	Unbind()
}

// Draw renders the quad.
func (d *SimulationDisplay) Draw(textures ...Bindable) {
	d.shader.Use()

	if d.transformDirty {
		m := d.transform.ComputeModel()
		m = m.Mul4(math.Scale3D(d.zoomFactor, d.zoomFactor, 1))
		d.shader.SetUniformMat4("Model", m)
		d.transformDirty = false
	}

	for i, tex := range textures {
		gl.ActiveTexture(gl.TEXTURE0 + uint32(i))
		tex.Bind()
	}

	gl.BindVertexArray(d.vao)
	gl.DrawArrays(gl.TRIANGLES, 0, 6)
	gl.BindVertexArray(0)

	for i, tex := range textures {
		gl.ActiveTexture(gl.TEXTURE0 + uint32(i))
		tex.Unbind()
	}

	d.shader.Unuse()
}
