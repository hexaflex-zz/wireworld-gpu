package main

import "strings"

// ShaderSource defines shader source code.
type ShaderSource struct {
	Vertex   string
	Geometry string
	Fragment string
}

// Compile compiles the given shader sources into a program.
func (s *ShaderSource) Compile() (Shader, error) {
	// Replace references to the shared source with the actual shared contents.
	const includeShared = "$INCLUDE_SHARED$"
	vs := strings.ReplaceAll(s.Vertex, includeShared, ShaderShared)
	gs := strings.ReplaceAll(s.Geometry, includeShared, ShaderShared)
	fs := strings.ReplaceAll(s.Fragment, includeShared, ShaderShared)
	return compile(string(vs), string(gs), string(fs))
}
