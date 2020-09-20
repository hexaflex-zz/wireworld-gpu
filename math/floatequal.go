// Copyright 2014 The go-gl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.go-gl file.

package math

// FloatEqual is a safe utility function to compare floats.
// It's Taken from http://floating-point-gui.de/errors/comparison/
//
// It is slightly altered to not call Abs when not needed.
func FloatEqual(a, b float32) bool {
	// Epsilon is some tiny value that determines how precisely equal we want our floats to be
	// This is exported and left as a variable in case you want to change the default threshold for the
	// purposes of certain methods (e.g. Unproject uses the default epsilon when determining
	// if the determinant is "close enough" to zero to mean there's no inverse).
	//
	// This is, obviously, not mutex protected so be **absolutely sure** that no functions using Epsilon
	// are being executed when you change this.
	const Epsilon = 1e-10

	return FloatEqualThreshold(a, b, Epsilon)
}

// FloatEqualFunc is a utility closure that will generate a function that
// always approximately compares floats like FloatEqualThreshold with a different
// threshold.
func FloatEqualFunc(epsilon float32) func(float32, float32) bool {
	return func(a, b float32) bool {
		return FloatEqualThreshold(a, b, epsilon)
	}
}

// FloatEqualThreshold is a utility function to compare floats.
// It's Taken from http://floating-point-gui.de/errors/comparison/
//
// It is slightly altered to not call Abs when not needed.
//
// This differs from FloatEqual in that it lets you pass in your comparison threshold, so that you can adjust the comparison value to your specific needs
func FloatEqualThreshold(a, b, epsilon float32) bool {
	// Various useful constants.
	var MinNormal = float32(1.1754943508222875e-38) // 1 / 2**(127 - 1)

	if a == b { // Handles the case of inf or shortcuts the loop when no significant error has accumulated
		return true
	}

	diff := Abs(a - b)
	if a*b == 0 || diff < MinNormal { // If a or b are 0 or both are extremely close to it
		return diff < epsilon*epsilon
	}

	// Else compare difference
	return diff/(Abs(a)+Abs(b)) < epsilon
}
