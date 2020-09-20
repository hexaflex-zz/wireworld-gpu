// Copyright 2014 The go-gl/mathgl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.go-gl file.

package math

import (
	"math"
)

// Vec3 defines a vector.
type Vec3 [3]float32

// DivScalar returns v / f.
func (v Vec3) DivScalar(f float32) Vec3 {
	return Vec3{
		v[0] / f,
		v[1] / f,
		v[2] / f,
	}
}

// Clamp returns v, with its components clamped to the range [min,max].
func (v Vec3) Clamp(min, max float32) Vec3 {
	v[0] = Clamp(v[0], min, max)
	v[1] = Clamp(v[1], min, max)
	v[2] = Clamp(v[2], min, max)
	return v
}

// Remainder returns a % b.
func (a Vec3) Remainder(b Vec3) Vec3 {
	return Vec3{
		Remainder(a[0], b[0]),
		Remainder(a[1], b[1]),
		Remainder(a[2], b[2]),
	}
}

// DivFloor returns floor(a / b).
func (a Vec3) DivFloor(b Vec3) Vec3 {
	return Vec3{
		Floor(a[0] / b[0]),
		Floor(a[1] / b[1]),
		Floor(a[2] / b[2]),
	}
}

// Div returns a / b.
func (a Vec3) Div(b Vec3) Vec3 {
	return Vec3{a[0] / b[0], a[1] / b[1], a[2] / b[2]}
}

// Mul returns a * b.
func (a Vec3) Mul(b Vec4) Vec3 {
	return Vec3{a[0] * b[0], a[1] * b[1], a[2] * b[2]}
}

// Add performs element-wise addition between two vectors. It is equivalent to iterating
// over every element of v1 and adding the corresponding element of v2 to it.
func (v1 Vec3) Add(v2 Vec3) Vec3 {
	return Vec3{v1[0] + v2[0], v1[1] + v2[1], v1[2] + v2[2]}
}

// Sub performs element-wise subtraction between two vectors. It is equivalent to iterating
// over every element of v1 and subtracting the corresponding element of v2 from it.
func (v1 Vec3) Sub(v2 Vec3) Vec3 {
	return Vec3{v1[0] - v2[0], v1[1] - v2[1], v1[2] - v2[2]}
}

// MulScalar performs a scalar multiplication between the vector and some constant value
// c. This is equivalent to iterating over every vector element and multiplying by c.
func (v1 Vec3) MulScalar(c float32) Vec3 {
	return Vec3{v1[0] * c, v1[1] * c, v1[2] * c}
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
func (v1 Vec3) Dot(v2 Vec3) float32 {
	return v1[0]*v2[0] + v1[1]*v2[1] + v1[2]*v2[2]
}

// Len returns the vector's length. Note that this is NOT the dimension of
// the vector (len(v)), but the mathematical length. This is equivalent to the square
// root of the sum of the squares of all elements. E.G. for a Vec2 it's
// math.Hypot(v[0], v[1]).
func (v1 Vec3) Len() float32 {

	return float32(math.Sqrt(float64(v1[0]*v1[0] + v1[1]*v1[1] + v1[2]*v1[2])))

}

// LenSqr returns the vector's square length. This is equivalent to the sum of the squares of all elements.
func (v1 Vec3) LenSqr() float32 {
	return v1[0]*v1[0] + v1[1]*v1[1] + v1[2]*v1[2]
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
func (v1 Vec3) Normalize() Vec3 {
	l := 1.0 / v1.Len()
	return Vec3{v1[0] * l, v1[1] * l, v1[2] * l}
}

// ApproxEqual takes in a vector and does an element-wise approximate float
// comparison as if FloatEqual had been used
func (v1 Vec3) ApproxEqual(v2 Vec3) bool {
	for i := range v1 {
		if !FloatEqual(v1[i], v2[i]) {
			return false
		}
	}
	return true
}

// ApproxEqualThreshold takes in a threshold for comparing two floats, and uses
// it to do an element-wise comparison of the vector to another.
func (v1 Vec3) ApproxEqualThreshold(v2 Vec3, threshold float32) bool {
	for i := range v1 {
		if !FloatEqualThreshold(v1[i], v2[i], threshold) {
			return false
		}
	}
	return true
}

// ApproxFuncEqual takes in a func that compares two floats, and uses it to do an element-wise
// comparison of the vector to another. This is intended to be used with FloatEqualFunc
func (v1 Vec3) ApproxFuncEqual(v2 Vec3, eq func(float32, float32) bool) bool {
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
func (v Vec3) X() float32 {
	return v[0]
}

// Y is an element access func, it is equivalent to v[n] where
// n is some valid index. The mappings are XYZW (X=0, Y=1 etc). Benchmarks
// show that this is more or less as fast as direct acces, probably due to
// inlining, so use v[0] or v.X() depending on personal preference.
func (v Vec3) Y() float32 {
	return v[1]
}

// Z is an element access func, it is equivalent to v[n] where
// n is some valid index. The mappings are XYZW (X=0, Y=1 etc). Benchmarks
// show that this is more or less as fast as direct acces, probably due to
// inlining, so use v[0] or v.X() depending on personal preference.
func (v Vec3) Z() float32 {
	return v[2]
}
