package main

import (
	"errors"
	"fmt"

	"github.com/go-gl/gl/v4.2-core/gl"
	math "github.com/hexaflex/glmath"
)

// SimulationState is an offscreen render target (framebuffer) which functions
// as the simulation state and applies the simulation rules.
type SimulationState struct {
	size math.Vec2
	fbo  uint32
	rbo  uint32
	tex  uint32
}

// Init initializes the framebuffer with the given size.
func (ss *SimulationState) Init(size math.Vec2) error {
	ss.size = size

	if size[0] < 1 || size[1] < 1 {
		return errors.New("framebuffer: invalid dimensions")
	}

	ss.Release()

	// Create framebuffer object.
	gl.GenFramebuffers(1, &ss.fbo)
	gl.BindFramebuffer(gl.FRAMEBUFFER, ss.fbo)

	// Create a texture object to store colour info.
	gl.GenTextures(1, &ss.tex)
	gl.BindTexture(gl.TEXTURE_2D, ss.tex)
	gl.TexParameterf(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexParameterf(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RED, int32(ss.size[0]), int32(ss.size[1]), 0, gl.RED, gl.UNSIGNED_BYTE, nil)
	gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0, gl.TEXTURE_2D, ss.tex, 0)

	// Create a render buffer object to store depth info.
	gl.GenRenderbuffers(1, &ss.rbo)
	gl.BindRenderbuffer(gl.RENDERBUFFER, ss.rbo)
	gl.RenderbufferStorage(gl.RENDERBUFFER, gl.DEPTH_COMPONENT, int32(ss.size[0]), int32(ss.size[1]))
	gl.FramebufferRenderbuffer(gl.FRAMEBUFFER, gl.DEPTH_ATTACHMENT, gl.RENDERBUFFER, ss.rbo)

	err := ss.checkStatus()

	gl.BindRenderbuffer(gl.RENDERBUFFER, 0)
	gl.BindTexture(gl.TEXTURE_2D, 0)
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
	return err
}

// Release clears framebuffer resources.
func (ss *SimulationState) Release() {
	gl.DeleteTextures(1, &ss.tex)
	gl.DeleteRenderbuffers(1, &ss.rbo)
	gl.DeleteFramebuffers(1, &ss.fbo)
}

// Size returns the framebuffer dimensions.
func (ss *SimulationState) Size() math.Vec2 {
	return ss.size
}

// BindTexture sets the framebuffer texture as the active texture.
// Call after rendering to the buffer is complete and you wish to
// use the buffer contents as a sampler in another drawing operation.
func (ss *SimulationState) BindTexture() {
	gl.BindTexture(gl.TEXTURE_2D, ss.tex)
}

// UnbindTexture unbinds the active texture.
func (ss *SimulationState) UnbindTexture() {
	gl.BindTexture(gl.TEXTURE_2D, 0)
}

// BindBuffer sets the buffer as the active render target.
func (ss *SimulationState) BindBuffer() {
	gl.BindFramebuffer(gl.FRAMEBUFFER, ss.fbo)
}

// UnbindBuffer unsets the buffer as the active render target.
func (ss *SimulationState) UnbindBuffer() {
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
}

// SetData writes the given state data into the framebuffer's color buffer.
// Sets the framebuffer dimensions to the given size.
func (ss *SimulationState) SetData(pix []byte, size math.Vec2) {
	ss.size = size

	gl.BindTexture(gl.TEXTURE_2D, ss.tex)
	gl.PixelStorei(gl.UNPACK_ALIGNMENT, 1)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RED, int32(ss.size[0]), int32(ss.size[1]), 0, gl.RED, gl.UNSIGNED_BYTE, gl.Ptr(pix))
	gl.BindTexture(gl.TEXTURE_2D, 0)
}

// Data reads state state from the framebuffer's color buffer.
// This uses glReadPixels and is therefore rather slow, so use with care.
func (ss *SimulationState) Data() []byte {
	p := make([]byte, int32(ss.size[0])*int32(ss.size[1]))
	gl.BindFramebuffer(gl.READ_FRAMEBUFFER, ss.fbo)
	gl.ReadPixels(0, 0, int32(ss.size[0]), int32(ss.size[1]), gl.RED, gl.UNSIGNED_BYTE, gl.Ptr(p))
	gl.BindFramebuffer(gl.READ_FRAMEBUFFER, 0)
	return p
}

func (ss *SimulationState) checkStatus() error {
	status := gl.CheckFramebufferStatus(gl.FRAMEBUFFER)

	switch status {
	case gl.FRAMEBUFFER_COMPLETE:
	case gl.FRAMEBUFFER_INCOMPLETE_ATTACHMENT:
		return errors.New("framebuffer error: Attachment is not complete")
	case gl.FRAMEBUFFER_INCOMPLETE_MISSING_ATTACHMENT:
		return errors.New("framebuffer error: no image attachment")
	case gl.FRAMEBUFFER_INCOMPLETE_DRAW_BUFFER:
		return errors.New("framebuffer error: draw buffer error")
	case gl.FRAMEBUFFER_INCOMPLETE_READ_BUFFER:
		return errors.New("framebuffer error: read buffer error")
	case gl.FRAMEBUFFER_INCOMPLETE_MULTISAMPLE:
		return errors.New("framebuffer error: incomplete multisample")
	case gl.FRAMEBUFFER_UNSUPPORTED:
		return errors.New("framebuffer error: FBO not supported")
	default:
		return fmt.Errorf("framebuffer error: unknown error %d", status)
	}
	return nil
}
