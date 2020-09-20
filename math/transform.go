package math

import (
	"math"
)

// Transform defines transformation state.
type Transform struct {
	Translate Vec2    // Translation offsets.
	Scale     Vec2    // Scale factors.
	Rotate    float32 // Rotation angle in radians.
}

// NewTransform returns a new, default transform.
func NewTransform() *Transform {
	return &Transform{
		Translate: Vec2{0, 0},
		Scale:     Vec2{1, 1},
		Rotate:    0,
	}
}

// ComputeModel computes and returns a model matrix from the transform state.
func (t *Transform) ComputeModel() Mat4 {
	translate := t.Translate
	rotate := t.Rotate
	scale := t.Scale

	model := Translate3D(translate[0]+0.375, translate[1]+0.375, 0)
	model = model.Mul4(HomogRotate3DZ(rotate))
	model = model.Mul4(Scale3D(scale[0], scale[1], 1))
	return model
}

// Copyright 2014 The go-gl/mathgl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.go-gl file.

// Translate3D returns a homogeneous (4x4 for 3D-space) Translation matrix that moves a point by Tx units in the x-direction, Ty units in the y-direction,
// and Tz units in the z-direction
//
//    [[1, 0, 0, Tx]]
//    [[0, 1, 0, Ty]]
//    [[0, 0, 1, Tz]]
//    [[0, 0, 0, 1 ]]
func Translate3D(Tx, Ty, Tz float32) Mat4 {
	return Mat4{1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, float32(Tx), float32(Ty), float32(Tz), 1}
}

// HomogRotate3DX is the same as Rotate3DX, except homogeneous (4x4 with the extra row/col being all zeroes with a one in the bottom right)
func HomogRotate3DX(angle float32) Mat4 {
	//angle = (angle * math.Pi) / 180.0
	sin, cos := float32(math.Sin(float64(angle))), float32(math.Cos(float64(angle)))

	return Mat4{1, 0, 0, 0, 0, cos, sin, 0, 0, -sin, cos, 0, 0, 0, 0, 1}
}

// HomogRotate3DY is the same as Rotate3DY, except homogeneous (4x4 with the extra row/col being all zeroes with a one in the bottom right)
func HomogRotate3DY(angle float32) Mat4 {
	//angle = (angle * math.Pi) / 180.0
	sin, cos := float32(math.Sin(float64(angle))), float32(math.Cos(float64(angle)))
	return Mat4{cos, 0, -sin, 0, 0, 1, 0, 0, sin, 0, cos, 0, 0, 0, 0, 1}
}

// HomogRotate3DZ is the same as Rotate3DZ, except homogeneous (4x4 with the extra row/col being all zeroes with a one in the bottom right)
func HomogRotate3DZ(angle float32) Mat4 {
	//angle = (angle * math.Pi) / 180.0
	sin, cos := float32(math.Sin(float64(angle))), float32(math.Cos(float64(angle)))
	return Mat4{cos, sin, 0, 0, -sin, cos, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1}
}

// Scale3D creates a homogeneous 3D scaling matrix
// [[ scaleX, 0     , 0     , 0 ]]
// [[ 0     , scaleY, 0     , 0 ]]
// [[ 0     , 0     , scaleZ, 0 ]]
// [[ 0     , 0     , 0     , 1 ]]
func Scale3D(scaleX, scaleY, scaleZ float32) Mat4 {

	return Mat4{float32(scaleX), 0, 0, 0, 0, float32(scaleY), 0, 0, 0, 0, float32(scaleZ), 0, 0, 0, 0, 1}
}

// ShearX3D creates a homogeneous 3D shear matrix along the X-axis
func ShearX3D(shearY, shearZ float32) Mat4 {
	return Mat4{1, float32(shearY), float32(shearZ), 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1}
}

// ShearY3D creates a homogeneous 3D shear matrix along the Y-axis
func ShearY3D(shearX, shearZ float32) Mat4 {
	return Mat4{1, 0, 0, 0, float32(shearX), 1, float32(shearZ), 0, 0, 0, 1, 0, 0, 0, 0, 1}
}

// ShearZ3D creates a homogeneous 3D shear matrix along the Z-axis
func ShearZ3D(shearX, shearY float32) Mat4 {
	return Mat4{1, 0, 0, 0, 0, 1, 0, 0, float32(shearX), float32(shearY), 1, 0, 0, 0, 0, 1}
}

// HomogRotate3D creates a 3D rotation Matrix that rotates by (radian) angle about some arbitrary axis given by a normalized Vector.
// It produces a homogeneous matrix (4x4)
//
// Where c is cos(angle) and s is sin(angle), and x, y, and z are the first, second, and third elements of the axis vector (respectively):
//
//    [[ x^2(1-c)+c, xy(1-c)-zs, xz(1-c)+ys, 0 ]]
//    [[ xy(1-c)+zs, y^2(1-c)+c, yz(1-c)-xs, 0 ]]
//    [[ xz(1-c)-ys, yz(1-c)+xs, z^2(1-c)+c, 0 ]]
//    [[ 0         , 0         , 0         , 1 ]]
func HomogRotate3D(angle float32, axis Vec3) Mat4 {
	x, y, z := axis[0], axis[1], axis[2]
	s, c := float32(math.Sin(float64(angle))), float32(math.Cos(float64(angle)))
	k := 1 - c

	return Mat4{x*x*k + c, x*y*k + z*s, x*z*k - y*s, 0, x*y*k - z*s, y*y*k + c, y*z*k + x*s, 0, x*z*k + y*s, y*z*k - x*s, z*z*k + c, 0, 0, 0, 0, 1}
}

// Extract3DScale extracts the 3d scaling from a homogeneous matrix
func Extract3DScale(m Mat4) (x, y, z float32) {
	return float32(math.Sqrt(float64(m[0]*m[0] + m[1]*m[1] + m[2]*m[2]))),
		float32(math.Sqrt(float64(m[4]*m[4] + m[5]*m[5] + m[6]*m[6]))),
		float32(math.Sqrt(float64(m[8]*m[8] + m[9]*m[9] + m[10]*m[10])))
}

// ExtractMaxScale extracts the maximum scaling from a homogeneous matrix
func ExtractMaxScale(m Mat4) float32 {
	scaleX := float64(m[0]*m[0] + m[1]*m[1] + m[2]*m[2])
	scaleY := float64(m[4]*m[4] + m[5]*m[5] + m[6]*m[6])
	scaleZ := float64(m[8]*m[8] + m[9]*m[9] + m[10]*m[10])

	return float32(math.Sqrt(math.Max(scaleX, math.Max(scaleY, scaleZ))))
}

// TransformCoordinate multiplies a 3D vector by a transformation given by
// the homogeneous 4D matrix m, applying any translation.
// If this transformation is non-affine, it will project this
// vector onto the plane w=1 before returning the result.
//
// This is similar to saying you're transforming and projecting a point.
//
// This is effectively equivalent to the GLSL code
//     vec4 r = (m * vec4(v,1.));
//     r = r/r.w;
//     vec3 newV = r.xyz;
func TransformCoordinate(v Vec3, m Mat4) Vec3 {
	t := v.Vec4(1)
	t = m.Mul4x1(t)
	t = t.MulScalar(1 / t[3])

	return t.Vec3()
}

// TransformNormal multiplies a 3D vector by a transformation given by
// the homogeneous 4D matrix m, NOT applying any translations.
//
// This is similar to saying you're applying a transformation
// to a direction or normal. Rotation still applies (as does scaling),
// but translating a direction or normal is meaningless.
//
// This is effectively equivalent to the GLSL code
//    vec4 r = (m * vec4(v,0.));
//    vec3 newV = r.xyz
func TransformNormal(v Vec3, m Mat4) Vec3 {
	t := v.Vec4(0)
	t = m.Mul4x1(t)

	return t.Vec3()
}
