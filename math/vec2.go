// Copyright 2014 The go-gl/mathgl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.go-gl file.

package math

import (
	"math"
)

// Vec2 defines a vector.
type Vec2 [2]float32

// DivScalar returns v / f.
func (v Vec2) DivScalar(f float32) Vec2 {
	return Vec2{
		v[0] / f,
		v[1] / f,
	}
}

// Clamp returns v, with its components clamped to the range [min,max].
func (v Vec2) Clamp(min, max float32) Vec2 {
	v[0] = Clamp(v[0], min, max)
	v[1] = Clamp(v[1], min, max)
	return v
}

// Remainder returns a % b.
func (a Vec2) Remainder(b Vec2) Vec2 {
	return Vec2{
		Remainder(a[0], b[0]),
		Remainder(a[1], b[1]),
	}
}

// DivFloor returns floor(a / b).
func (a Vec2) DivFloor(b Vec2) Vec2 {
	return Vec2{
		Floor(a[0] / b[0]),
		Floor(a[1] / b[1]),
	}
}

// Div returns a / b.
func (a Vec2) Div(b Vec2) Vec2 {
	return Vec2{a[0] / b[0], a[1] / b[1]}
}

// Mul returns a * b.
func (a Vec2) Mul(b Vec2) Vec2 {
	return Vec2{a[0] * b[0], a[1] * b[1]}
}

// Vec3 constructs a 3-dimensional vector by appending the given coordinates.
func (v Vec2) Vec3(z float32) Vec3 {
	return Vec3{v[0], v[1], z}
}

// Vec4 constructs a 4-dimensional vector by appending the given coordinates.
func (v Vec2) Vec4(z, w float32) Vec4 {
	return Vec4{v[0], v[1], z, w}
}

// Vec4 constructs a 4-dimensional vector by appending the given coordinates.
func (v Vec3) Vec4(w float32) Vec4 {
	return Vec4{v[0], v[1], v[2], w}
}

// Vec2 constructs a 2-dimensional vector by discarding coordinates.
func (v Vec3) Vec2() Vec2 {
	return Vec2{v[0], v[1]}
}

// Vec2 constructs a 2-dimensional vector by discarding coordinates.
func (v Vec4) Vec2() Vec2 {
	return Vec2{v[0], v[1]}
}

// Vec3 constructs a 3-dimensional vector by discarding coordinates.
func (v Vec4) Vec3() Vec3 {
	return Vec3{v[0], v[1], v[2]}
}

// Elem extracts the elements of the vector for direct value assignment.
func (v Vec2) Elem() (x, y float32) {
	return v[0], v[1]
}

// Elem extracts the elements of the vector for direct value assignment.
func (v Vec3) Elem() (x, y, z float32) {
	return v[0], v[1], v[2]
}

// Elem extracts the elements of the vector for direct value assignment.
func (v Vec4) Elem() (x, y, z, w float32) {
	return v[0], v[1], v[2], v[3]
}

// Cross is the vector cross product. This operation is only defined on 3D
// vectors. It is equivalent to Vec3{v1[1]*v2[2]-v1[2]*v2[1],
// v1[2]*v2[0]-v1[0]*v2[2], v1[0]*v2[1] - v1[1]*v2[0]}. Another interpretation
// is that it's the vector whose magnitude is |v1||v2|sin(theta) where theta is
// the angle between v1 and v2.
//
// The cross product is most often used for finding surface normals. The cross
// product of vectors will generate a vector that is perpendicular to the plane
// they form.
//
// Technically, a generalized cross product exists as an "(N-1)ary" operation
// (that is, the 4D cross product requires 3 4D vectors). But the binary 3D (and
// 7D) cross product is the most important. It can be considered the area of a
// parallelogram with sides v1 and v2.
//
// Like the dot product, the cross product is roughly a measure of
// directionality. Two normalized perpendicular vectors will return a vector
// with a magnitude of 1.0 or -1.0 and two parallel vectors will return a vector
// with magnitude 0.0. The cross product is "anticommutative" meaning
// v1.Cross(v2) = -v2.Cross(v1), this property can be useful to know when
// finding normals, as taking the wrong cross product can lead to the opposite
// normal of the one you want.
func (v1 Vec3) Cross(v2 Vec3) Vec3 {
	return Vec3{v1[1]*v2[2] - v1[2]*v2[1], v1[2]*v2[0] - v1[0]*v2[2], v1[0]*v2[1] - v1[1]*v2[0]}
}

// Add performs element-wise addition between two vectors. It is equivalent to iterating
// over every element of v1 and adding the corresponding element of v2 to it.
func (v1 Vec2) Add(v2 Vec2) Vec2 {
	return Vec2{v1[0] + v2[0], v1[1] + v2[1]}
}

// Sub performs element-wise subtraction between two vectors. It is equivalent to iterating
// over every element of v1 and subtracting the corresponding element of v2 from it.
func (v1 Vec2) Sub(v2 Vec2) Vec2 {
	return Vec2{v1[0] - v2[0], v1[1] - v2[1]}
}

// MulScalar performs a scalar multiplication between the vector and some constant value
// c. This is equivalent to iterating over every vector element and multiplying by c.
func (v1 Vec2) MulScalar(c float32) Vec2 {
	return Vec2{v1[0] * c, v1[1] * c}
}

// Dot returns the dot product of this vector with another. There are multiple ways
// to describe this value. One is the multiplication of their lengths and cos(theta) where
// theta is the angle between the vectors: v1.v2 = |v1||v2|cos(theta).
//
// The other (and what is actually done) is the sum of the element-wise multiplication of all
// elements. So for instance, two Vec3s would yield v1.x * v2.x + v1.y * v2.y + v1.z * v2.z.
//
// This means that the dot product of a vector and itself is the square of its Len (within
// the bounds of floating points error).
//
// The dot product is roughly a measure of how closely two vectors are to pointing in the same
// direction. If both vectors are normalized, the value will be -1 for opposite pointing,
// one for same pointing, and 0 for perpendicular vectors.
func (v1 Vec2) Dot(v2 Vec2) float32 {
	return v1[0]*v2[0] + v1[1]*v2[1]
}

// Len returns the vector's length. Note that this is NOT the dimension of
// the vector (len(v)), but the mathematical length. This is equivalent to the square
// root of the sum of the squares of all elements. E.G. for a Vec2 it's
// math.Hypot(v[0], v[1]).
func (v1 Vec2) Len() float32 {

	return float32(math.Hypot(float64(v1[0]), float64(v1[1])))

}

// LenSqr returns the vector's square length. This is equivalent to the sum of the squares of all elements.
func (v1 Vec2) LenSqr() float32 {
	return v1[0]*v1[0] + v1[1]*v1[1]
}

// Normalize normalizes the vector. Normalization is (1/|v|)*v,
// making this equivalent to v.Scale(1/v.Len()). If the len is 0.0,
// this function will return an infinite value for all elements due
// to how floating point division works in Go (n/0.0 = math.Inf(Sign(n))).
//
// Normalization makes a vector's Len become 1.0 (within the margin of floating point error),
// while maintaining its directionality.
//
// (Can be seen here: http://play.golang.org/p/Aaj7SnbqIp )
func (v1 Vec2) Normalize() Vec2 {
	l := 1.0 / v1.Len()
	return Vec2{v1[0] * l, v1[1] * l}
}

// ApproxEqual takes in a vector and does an element-wise approximate float
// comparison as if FloatEqual had been used
func (v1 Vec2) ApproxEqual(v2 Vec2) bool {
	for i := range v1 {
		if !FloatEqual(v1[i], v2[i]) {
			return false
		}
	}
	return true
}

// ApproxEqualThreshold takes in a threshold for comparing two floats, and uses
// it to do an element-wise comparison of the vector to another.
func (v1 Vec2) ApproxEqualThreshold(v2 Vec2, threshold float32) bool {
	for i := range v1 {
		if !FloatEqualThreshold(v1[i], v2[i], threshold) {
			return false
		}
	}
	return true
}

// ApproxFuncEqual takes in a func that compares two floats, and uses it to do an element-wise
// comparison of the vector to another. This is intended to be used with FloatEqualFunc
func (v1 Vec2) ApproxFuncEqual(v2 Vec2, eq func(float32, float32) bool) bool {
	for i := range v1 {
		if !eq(v1[i], v2[i]) {
			return false
		}
	}
	return true
}

// X is an element access func, it is equivalent to v[n] where
// n is some valid index. The mappings are XYZW (X=0, Y=1 etc). Benchmarks
// show that this is more or less as fast as direct acces, probably due to
// inlining, so use v[0] or v.X() depending on personal preference.
func (v Vec2) X() float32 {
	return v[0]
}

// Y is an element access func, it is equivalent to v[n] where
// n is some valid index. The mappings are XYZW (X=0, Y=1 etc). Benchmarks
// show that this is more or less as fast as direct acces, probably due to
// inlining, so use v[0] or v.X() depending on personal preference.
func (v Vec2) Y() float32 {
	return v[1]
}
