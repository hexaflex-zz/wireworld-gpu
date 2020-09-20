// Copyright 2014 The go-gl/mathgl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.go-gl file.

package math

import (
	"bytes"
	"fmt"
	"text/tabwriter"
)

// Mat4 is a 4x4 matrix in row major order.
//
// m[4*r + c] is the element in the r'th row and c'th column.
type Mat4 [16]float32

// SetCol sets a Column within the Matrix, so it mutates the calling matrix.
func (m1 *Mat4) SetCol(col int, v Vec4) {
	m1[col*4+0], m1[col*4+1], m1[col*4+2], m1[col*4+3] = v[0], v[1], v[2], v[3]
}

// SetRow sets a Row within the Matrix, so it mutates the calling matrix.
func (m1 *Mat4) SetRow(row int, v Vec4) {
	m1[row+0], m1[row+4], m1[row+8], m1[row+12] = v[0], v[1], v[2], v[3]
}

// Diag is a basic operation on a square matrix that simply
// returns main diagonal (meaning all elements such that row==col).
func (m1 Mat4) Diag() Vec4 {
	return Vec4{m1[0], m1[5], m1[10], m1[15]}
}

// Ident4 returns the 4x4 identity matrix.
// The identity matrix is a square matrix with the value 1 on its
// diagonals. The characteristic property of the identity matrix is that
// any matrix multiplied by it is itself. (MI = M; IN = N)
func Ident4() Mat4 {
	return Mat4{1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1}
}

// Add performs an element-wise addition of two matrices, this is
// equivalent to iterating over every element of m1 and adding the corresponding value of m2.
func (m1 Mat4) Add(m2 Mat4) Mat4 {
	return Mat4{m1[0] + m2[0], m1[1] + m2[1], m1[2] + m2[2], m1[3] + m2[3], m1[4] + m2[4], m1[5] + m2[5], m1[6] + m2[6], m1[7] + m2[7], m1[8] + m2[8], m1[9] + m2[9], m1[10] + m2[10], m1[11] + m2[11], m1[12] + m2[12], m1[13] + m2[13], m1[14] + m2[14], m1[15] + m2[15]}
}

// Sub performs an element-wise subtraction of two matrices, this is
// equivalent to iterating over every element of m1 and subtracting the corresponding value of m2.
func (m1 Mat4) Sub(m2 Mat4) Mat4 {
	return Mat4{m1[0] - m2[0], m1[1] - m2[1], m1[2] - m2[2], m1[3] - m2[3], m1[4] - m2[4], m1[5] - m2[5], m1[6] - m2[6], m1[7] - m2[7], m1[8] - m2[8], m1[9] - m2[9], m1[10] - m2[10], m1[11] - m2[11], m1[12] - m2[12], m1[13] - m2[13], m1[14] - m2[14], m1[15] - m2[15]}
}

// Mul performs a scalar multiplcation of the matrix. This is equivalent to iterating
// over every element of the matrix and multiply it by c.
func (m1 Mat4) Mul(c float32) Mat4 {
	return Mat4{m1[0] * c, m1[1] * c, m1[2] * c, m1[3] * c, m1[4] * c, m1[5] * c, m1[6] * c, m1[7] * c, m1[8] * c, m1[9] * c, m1[10] * c, m1[11] * c, m1[12] * c, m1[13] * c, m1[14] * c, m1[15] * c}
}

// Mul4x1 performs a "matrix product" between this matrix
// and another of the given dimension. For any two matrices of dimensionality
// MxN and NxO, the result will be MxO. For instance, Mat4 multiplied using
// Mul4x2 will result in a Mat4x2.
func (m1 Mat4) Mul4x1(m2 Vec4) Vec4 {
	return Vec4{
		m1[0]*m2[0] + m1[4]*m2[1] + m1[8]*m2[2] + m1[12]*m2[3],
		m1[1]*m2[0] + m1[5]*m2[1] + m1[9]*m2[2] + m1[13]*m2[3],
		m1[2]*m2[0] + m1[6]*m2[1] + m1[10]*m2[2] + m1[14]*m2[3],
		m1[3]*m2[0] + m1[7]*m2[1] + m1[11]*m2[2] + m1[15]*m2[3],
	}
}

// Mul4 performs a "matrix product" between this matrix
// and another of the given dimension. For any two matrices of dimensionality
// MxN and NxO, the result will be MxO. For instance, Mat4 multiplied using
// Mul4x2 will result in a Mat4x2.
func (m1 Mat4) Mul4(m2 Mat4) Mat4 {
	return Mat4{
		m1[0]*m2[0] + m1[4]*m2[1] + m1[8]*m2[2] + m1[12]*m2[3],
		m1[1]*m2[0] + m1[5]*m2[1] + m1[9]*m2[2] + m1[13]*m2[3],
		m1[2]*m2[0] + m1[6]*m2[1] + m1[10]*m2[2] + m1[14]*m2[3],
		m1[3]*m2[0] + m1[7]*m2[1] + m1[11]*m2[2] + m1[15]*m2[3],
		m1[0]*m2[4] + m1[4]*m2[5] + m1[8]*m2[6] + m1[12]*m2[7],
		m1[1]*m2[4] + m1[5]*m2[5] + m1[9]*m2[6] + m1[13]*m2[7],
		m1[2]*m2[4] + m1[6]*m2[5] + m1[10]*m2[6] + m1[14]*m2[7],
		m1[3]*m2[4] + m1[7]*m2[5] + m1[11]*m2[6] + m1[15]*m2[7],
		m1[0]*m2[8] + m1[4]*m2[9] + m1[8]*m2[10] + m1[12]*m2[11],
		m1[1]*m2[8] + m1[5]*m2[9] + m1[9]*m2[10] + m1[13]*m2[11],
		m1[2]*m2[8] + m1[6]*m2[9] + m1[10]*m2[10] + m1[14]*m2[11],
		m1[3]*m2[8] + m1[7]*m2[9] + m1[11]*m2[10] + m1[15]*m2[11],
		m1[0]*m2[12] + m1[4]*m2[13] + m1[8]*m2[14] + m1[12]*m2[15],
		m1[1]*m2[12] + m1[5]*m2[13] + m1[9]*m2[14] + m1[13]*m2[15],
		m1[2]*m2[12] + m1[6]*m2[13] + m1[10]*m2[14] + m1[14]*m2[15],
		m1[3]*m2[12] + m1[7]*m2[13] + m1[11]*m2[14] + m1[15]*m2[15],
	}
}

// Transpose produces the transpose of this matrix. For any MxN matrix
// the transpose is an NxM matrix with the rows swapped with the columns. For instance
// the transpose of the Mat3x2 is a Mat2x3 like so:
//
//    [[a b]]    [[a c e]]
//    [[c d]] =  [[b d f]]
//    [[e f]]
func (m1 Mat4) Transpose() Mat4 {
	return Mat4{m1[0], m1[4], m1[8], m1[12], m1[1], m1[5], m1[9], m1[13], m1[2], m1[6], m1[10], m1[14], m1[3], m1[7], m1[11], m1[15]}
}

// Det returns the determinant of a matrix. It is a measure of a square matrix's
// singularity and invertability, among other things. In this library, the
// determinant is hard coded based on pre-computed cofactor expansion, and uses
// no loops. Of course, the addition and multiplication must still be done.
func (m1 Mat4) Det() float32 {
	return m1[0]*m1[5]*m1[10]*m1[15] - m1[0]*m1[5]*m1[11]*m1[14] - m1[0]*m1[6]*m1[9]*m1[15] + m1[0]*m1[6]*m1[11]*m1[13] + m1[0]*m1[7]*m1[9]*m1[14] - m1[0]*m1[7]*m1[10]*m1[13] - m1[1]*m1[4]*m1[10]*m1[15] + m1[1]*m1[4]*m1[11]*m1[14] + m1[1]*m1[6]*m1[8]*m1[15] - m1[1]*m1[6]*m1[11]*m1[12] - m1[1]*m1[7]*m1[8]*m1[14] + m1[1]*m1[7]*m1[10]*m1[12] + m1[2]*m1[4]*m1[9]*m1[15] - m1[2]*m1[4]*m1[11]*m1[13] - m1[2]*m1[5]*m1[8]*m1[15] + m1[2]*m1[5]*m1[11]*m1[12] + m1[2]*m1[7]*m1[8]*m1[13] - m1[2]*m1[7]*m1[9]*m1[12] - m1[3]*m1[4]*m1[9]*m1[14] + m1[3]*m1[4]*m1[10]*m1[13] + m1[3]*m1[5]*m1[8]*m1[14] - m1[3]*m1[5]*m1[10]*m1[12] - m1[3]*m1[6]*m1[8]*m1[13] + m1[3]*m1[6]*m1[9]*m1[12]
}

// Inv computes the inverse of a square matrix. An inverse is a square matrix such that when multiplied by the
// original, yields the identity.
//
// M_inv * M = M * M_inv = I
//
// In this library, the math is precomputed, and uses no loops, though the multiplications, additions, determinant calculation, and scaling
// are still done. This can still be (relatively) expensive for a 4x4.
//
// This function checks the determinant to see if the matrix is invertible.
// If the determinant is 0.0, this function returns the zero matrix. However, due to floating point errors, it is
// entirely plausible to get a false positive or negative.
// In the future, an alternate function may be written which takes in a pre-computed determinant.
func (m1 Mat4) Inv() Mat4 {
	det := m1.Det()
	if FloatEqual(det, float32(0.0)) {
		return Mat4{}
	}

	retMat := Mat4{
		-m1[7]*m1[10]*m1[13] + m1[6]*m1[11]*m1[13] + m1[7]*m1[9]*m1[14] - m1[5]*m1[11]*m1[14] - m1[6]*m1[9]*m1[15] + m1[5]*m1[10]*m1[15],
		m1[3]*m1[10]*m1[13] - m1[2]*m1[11]*m1[13] - m1[3]*m1[9]*m1[14] + m1[1]*m1[11]*m1[14] + m1[2]*m1[9]*m1[15] - m1[1]*m1[10]*m1[15],
		-m1[3]*m1[6]*m1[13] + m1[2]*m1[7]*m1[13] + m1[3]*m1[5]*m1[14] - m1[1]*m1[7]*m1[14] - m1[2]*m1[5]*m1[15] + m1[1]*m1[6]*m1[15],
		m1[3]*m1[6]*m1[9] - m1[2]*m1[7]*m1[9] - m1[3]*m1[5]*m1[10] + m1[1]*m1[7]*m1[10] + m1[2]*m1[5]*m1[11] - m1[1]*m1[6]*m1[11],
		m1[7]*m1[10]*m1[12] - m1[6]*m1[11]*m1[12] - m1[7]*m1[8]*m1[14] + m1[4]*m1[11]*m1[14] + m1[6]*m1[8]*m1[15] - m1[4]*m1[10]*m1[15],
		-m1[3]*m1[10]*m1[12] + m1[2]*m1[11]*m1[12] + m1[3]*m1[8]*m1[14] - m1[0]*m1[11]*m1[14] - m1[2]*m1[8]*m1[15] + m1[0]*m1[10]*m1[15],
		m1[3]*m1[6]*m1[12] - m1[2]*m1[7]*m1[12] - m1[3]*m1[4]*m1[14] + m1[0]*m1[7]*m1[14] + m1[2]*m1[4]*m1[15] - m1[0]*m1[6]*m1[15],
		-m1[3]*m1[6]*m1[8] + m1[2]*m1[7]*m1[8] + m1[3]*m1[4]*m1[10] - m1[0]*m1[7]*m1[10] - m1[2]*m1[4]*m1[11] + m1[0]*m1[6]*m1[11],
		-m1[7]*m1[9]*m1[12] + m1[5]*m1[11]*m1[12] + m1[7]*m1[8]*m1[13] - m1[4]*m1[11]*m1[13] - m1[5]*m1[8]*m1[15] + m1[4]*m1[9]*m1[15],
		m1[3]*m1[9]*m1[12] - m1[1]*m1[11]*m1[12] - m1[3]*m1[8]*m1[13] + m1[0]*m1[11]*m1[13] + m1[1]*m1[8]*m1[15] - m1[0]*m1[9]*m1[15],
		-m1[3]*m1[5]*m1[12] + m1[1]*m1[7]*m1[12] + m1[3]*m1[4]*m1[13] - m1[0]*m1[7]*m1[13] - m1[1]*m1[4]*m1[15] + m1[0]*m1[5]*m1[15],
		m1[3]*m1[5]*m1[8] - m1[1]*m1[7]*m1[8] - m1[3]*m1[4]*m1[9] + m1[0]*m1[7]*m1[9] + m1[1]*m1[4]*m1[11] - m1[0]*m1[5]*m1[11],
		m1[6]*m1[9]*m1[12] - m1[5]*m1[10]*m1[12] - m1[6]*m1[8]*m1[13] + m1[4]*m1[10]*m1[13] + m1[5]*m1[8]*m1[14] - m1[4]*m1[9]*m1[14],
		-m1[2]*m1[9]*m1[12] + m1[1]*m1[10]*m1[12] + m1[2]*m1[8]*m1[13] - m1[0]*m1[10]*m1[13] - m1[1]*m1[8]*m1[14] + m1[0]*m1[9]*m1[14],
		m1[2]*m1[5]*m1[12] - m1[1]*m1[6]*m1[12] - m1[2]*m1[4]*m1[13] + m1[0]*m1[6]*m1[13] + m1[1]*m1[4]*m1[14] - m1[0]*m1[5]*m1[14],
		-m1[2]*m1[5]*m1[8] + m1[1]*m1[6]*m1[8] + m1[2]*m1[4]*m1[9] - m1[0]*m1[6]*m1[9] - m1[1]*m1[4]*m1[10] + m1[0]*m1[5]*m1[10],
	}

	return retMat.Mul(1 / det)
}

// ApproxEqual performs an element-wise approximate equality test between two matrices,
// as if FloatEqual had been used.
func (m1 Mat4) ApproxEqual(m2 Mat4) bool {
	for i := range m1 {
		if !FloatEqual(m1[i], m2[i]) {
			return false
		}
	}
	return true
}

// ApproxEqualThreshold performs an element-wise approximate equality test between two matrices
// with a given epsilon threshold, as if FloatEqualThreshold had been used.
func (m1 Mat4) ApproxEqualThreshold(m2 Mat4, threshold float32) bool {
	for i := range m1 {
		if !FloatEqualThreshold(m1[i], m2[i], threshold) {
			return false
		}
	}
	return true
}

// ApproxFuncEqual performs an element-wise approximate equality test between two matrices
// with a given equality functions, intended to be used with FloatEqualFunc; although and comparison
// function may be used in practice.
func (m1 Mat4) ApproxFuncEqual(m2 Mat4, eq func(float32, float32) bool) bool {
	for i := range m1 {
		if !eq(m1[i], m2[i]) {
			return false
		}
	}
	return true
}

// At returns the matrix element at the given row and column.
// This is equivalent to mat[col * numRow + row] where numRow is constant
// (E.G. for a Mat3x2 it's equal to 3)
//
// This method is garbage-in garbage-out. For instance, on a Mat4 asking for
// At(5,0) will work just like At(1,1). Or it may panic if it's out of bounds.
func (m1 Mat4) At(row, col int) float32 {
	return m1[col*4+row]
}

// Set sets the corresponding matrix element at the given row and column.
// This has a pointer receiver because it mutates the matrix.
//
// This method is garbage-in garbage-out. For instance, on a Mat4 asking for
// Set(5,0,val) will work just like Set(1,1,val). Or it may panic if it's out of bounds.
func (m1 *Mat4) Set(row, col int, value float32) {
	m1[col*4+row] = value
}

// Index returns the index of the given row and column, to be used with direct
// access. E.G. Index(0,0) = 0.
//
// This is a garbage-in garbage-out method. For instance, on a Mat4 asking for the index of
// (5,0) will work the same as asking for (1,1). Or it may give you a value that will cause
// a panic if you try to access the array with it if it's truly out of bounds.
func (m1 Mat4) Index(row, col int) int {
	return col*4 + row
}

// Row returns a vector representing the corresponding row (starting at row 0).
// This package makes no distinction between row and column vectors, so it
// will be a normal VecM for a MxN matrix.
func (m1 Mat4) Row(row int) Vec4 {
	return Vec4{m1[row+0], m1[row+4], m1[row+8], m1[row+12]}
}

// Rows decomposes a matrix into its corresponding row vectors.
// This is equivalent to calling mat.Row for each row.
func (m1 Mat4) Rows() (row0, row1, row2, row3 Vec4) {
	return m1.Row(0), m1.Row(1), m1.Row(2), m1.Row(3)
}

// Col returns a vector representing the corresponding column (starting at col 0).
// This package makes no distinction between row and column vectors, so it
// will be a normal VecN for a MxN matrix.
func (m1 Mat4) Col(col int) Vec4 {
	return Vec4{m1[col*4+0], m1[col*4+1], m1[col*4+2], m1[col*4+3]}
}

// Cols decomposes a matrix into its corresponding column vectors.
// This is equivalent to calling mat.Col for each column.
func (m1 Mat4) Cols() (col0, col1, col2, col3 Vec4) {
	return m1.Col(0), m1.Col(1), m1.Col(2), m1.Col(3)
}

// Trace is a basic operation on a square matrix that simply
// sums up all elements on the main diagonal (meaning all elements such that row==col).
func (m1 Mat4) Trace() float32 {
	return m1[0] + m1[5] + m1[10] + m1[15]
}

// Abs returns the element-wise absolute value of this matrix
func (m1 Mat4) Abs() Mat4 {
	return Mat4{Abs(m1[0]), Abs(m1[1]), Abs(m1[2]), Abs(m1[3]), Abs(m1[4]), Abs(m1[5]), Abs(m1[6]), Abs(m1[7]), Abs(m1[8]), Abs(m1[9]), Abs(m1[10]), Abs(m1[11]), Abs(m1[12]), Abs(m1[13]), Abs(m1[14]), Abs(m1[15])}
}

// Pretty prints the matrix
func (m1 Mat4) String() string {
	buf := new(bytes.Buffer)
	w := tabwriter.NewWriter(buf, 4, 4, 1, ' ', tabwriter.AlignRight)
	for i := 0; i < 4; i++ {
		for _, col := range m1.Row(i) {
			fmt.Fprintf(w, "%f\t", col)
		}

		fmt.Fprintln(w, "")
	}
	w.Flush()

	return buf.String()
}
