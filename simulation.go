package main

import (
	"image"
	"os"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	_ "github.com/hexaflex/pnm"

	"github.com/go-gl/gl/v4.2-core/gl"
	math "github.com/hexaflex/glmath"
)

// Simulation implements the GPU driven wireworld simulation.
type Simulation struct {
	shader Shader
	input  SimulationState
	output SimulationState
	vao    uint32
	vbo    uint32
}

// NewSimulation creates a new, empty simulation with the given dimensions.
func NewSimulation(width, height int) (*Simulation, error) {
	var err error
	var s Simulation

	size := math.Vec2{float32(width), float32(height)}

	s.shader, err = SimulationShader.Compile()
	if err != nil {
		return nil, err
	}

	if err = s.input.Init(size); err != nil {
		return nil, err
	}

	if err = s.output.Init(size); err != nil {
		s.Release()
		return nil, err
	}

	var verts = []float32{
		// x,y,u,v
		-1, -1, 0, 0,
		1, -1, 1, 0,
		-1, 1, 0, 1,
		1, -1, 1, 0,
		1, 1, 1, 1,
		-1, 1, 0, 1}

	gl.GenVertexArrays(1, &s.vao)
	gl.BindVertexArray(s.vao)

	gl.GenBuffers(1, &s.vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, s.vbo)
	gl.EnableVertexAttribArray(0)
	gl.EnableVertexAttribArray(1)
	gl.VertexAttribPointer(0, 2, gl.FLOAT, false, 4*4, gl.PtrOffset(0))
	gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 4*4, gl.PtrOffset(2*4))
	gl.BufferData(gl.ARRAY_BUFFER, len(verts)*4, gl.Ptr(verts), gl.STATIC_DRAW)

	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindVertexArray(0)

	return &s, nil
}

// LoadSimulation loads a simulation from the given image file.
// Supported formats: PNG, JPG, GIF, PNM
//
// It uses the given color palette to recognize cell states.
func LoadSimulation(file string, pal *Palette) (*Simulation, error) {
	fd, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	img, _, err := image.Decode(fd)
	fd.Close()
	if err != nil {
		return nil, err
	}

	b := img.Bounds()
	sim, err := NewSimulation(b.Dx(), b.Dy())
	if err != nil {
		return nil, err
	}

	// Set the input buffer to the image data.
	pix, size := pal.ToInternalFormat(img)
	sim.input.SetData(pix, size)
	return sim, nil
}

// Release unloads simulator resources.
func (s *Simulation) Release() {
	gl.DeleteBuffers(1, &s.vbo)
	gl.DeleteVertexArrays(1, &s.vao)
	s.shader.Release()
	s.input.Release()
	s.output.Release()
}

// Size returns the cell dimensions of the simulation.
func (s *Simulation) Size() math.Vec2 {
	return s.output.Size()
}

// Image returns the current simulation state as an image,
// colored using the given palette. Note that this uses
// glReadPixels and consequently is rather slow. Use it sparingly.
func (s *Simulation) Image(pal *Palette) image.Image {
	// We read from input because the render function sets
	// this to the most recent simulation state.
	size := s.input.Size()
	pix := s.input.Data()
	return pal.fromInternalFormat(pix, size)
}

// Bind binds the current simulation state's texture, so it may be
// used in other rendering operations.
func (s *Simulation) Bind() {
	s.input.BindTexture()
}

// Unbind unbinds the current texture.
func (s *Simulation) Unbind() {
	s.input.UnbindTexture()
}

// Step runs the simulation n times.
func (s *Simulation) Step(n int) {
	if n < 1 {
		return
	}
	s.shader.Use()

	size := s.input.Size()
	gl.Viewport(0, 0, int32(size[0]), int32(size[1]))
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.BindVertexArray(s.vao)
	gl.ActiveTexture(gl.TEXTURE0)

	for i := 0; i < n; i++ {
		s.output.BindBuffer()
		s.input.BindTexture()

		gl.DrawArrays(gl.TRIANGLES, 0, 6)

		s.input.UnbindTexture()
		s.output.UnbindBuffer()

		// Swap the states around. So the output of this pass
		// becomes the input of the next pass.
		s.output, s.input = s.input, s.output
	}

	gl.BindVertexArray(0)
	s.shader.Unuse()
}
