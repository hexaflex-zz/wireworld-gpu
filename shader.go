package main

import (
	"fmt"
	"strings"

	"github.com/go-gl/gl/v4.2-core/gl"
	"github.com/hexaflex/wireworld-gpu/math"
	"github.com/pkg/errors"
)

// Shader defines a compiled shader program.
type Shader uint32

// Release cleans up the shader.
func (s Shader) Release() {
	gl.DeleteShader(uint32(s))
}

// SetUniformMat4 sets the given uniform to the specified value.
func (s Shader) SetUniformMat4(name string, mat math.Mat4) {
	gl.UniformMatrix4fv(s.uniform(name), 1, false, &mat[0])
}

// SetUniformVec2 sets the given uniform to the specified value.
func (s Shader) SetUniformVec2(name string, v math.Vec2) {
	gl.Uniform2fv(s.uniform(name), 1, &v[0])
}

// SetUniformVec3 sets the given uniform to the specified value.
func (s Shader) SetUniformVec3(name string, v math.Vec3) {
	gl.Uniform3fv(s.uniform(name), 1, &v[0])
}

// SetUniformVec4 sets the given uniform to the specified value.
func (s Shader) SetUniformVec4(name string, v math.Vec4) {
	gl.Uniform4fv(s.uniform(name), 1, &v[0])
}

func (s Shader) uniform(name string) int32 {
	return gl.GetUniformLocation(uint32(s), gl.Str(name+"\x00"))
}

// Use uses the Shader.
func (s Shader) Use() {
	gl.UseProgram(uint32(s))
}

// Unuse unuses the Shader.
func (s Shader) Unuse() {
	gl.UseProgram(0)
}

// compile loads a shader from the given sources.
func compile(vertex, geometry, fragment string) (Shader, error) {
	var vs, gs, fs uint32
	var err error

	if len(vertex) > 0 {
		vs, err = compileShader(vertex, gl.VERTEX_SHADER)
		if err != nil {
			return 0, errors.Wrapf(err, "failed to compile vertex shader")
		}

		defer gl.DeleteShader(vs)
	}

	if len(geometry) > 0 {
		gs, err = compileShader(geometry, gl.GEOMETRY_SHADER)
		if err != nil {
			return 0, errors.Wrapf(err, "failed to compile geometry shader")
		}

		defer gl.DeleteShader(gs)
	}

	if len(fragment) > 0 {
		fs, err = compileShader(fragment, gl.FRAGMENT_SHADER)
		if err != nil {
			return 0, errors.Wrapf(err, "failed to compile fragment shader")
		}

		defer gl.DeleteShader(fs)
	}

	program := gl.CreateProgram()

	if len(vertex) > 0 {
		gl.AttachShader(program, vs)
	}

	if len(geometry) > 0 {
		gl.AttachShader(program, gs)
	}

	if len(fragment) > 0 {
		gl.AttachShader(program, fs)
	}

	gl.LinkProgram(program)

	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(program, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to link program: %v", log)
	}

	return Shader(program), nil
}

// compileShader compiles the given shader source into a Shader.
func compileShader(source string, stype uint32) (uint32, error) {
	shader := gl.CreateShader(stype)

	csources, free := gl.Strs(source + "\x00")
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}

	return shader, nil
}
