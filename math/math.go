package math

import "math"

// Pow2 returns the first power-of-two value >= to n.
// This can be used to create suitable texture dimensions.
func Pow2(n int) int {
	x := uint32(n) - 1
	x |= x >> 1
	x |= x >> 2
	x |= x >> 4
	x |= x >> 8
	x |= x >> 16
	return int(x + 1)
}

// IsPow2 returns true if the given value is a power-of-two.
func IsPow2(n int) bool {
	x := uint32(n)
	return (x & (x - 1)) == 0
}

// Clamp returns x, clamped to the range [min,max].
func Clamp(x, min, max float32) float32 {
	if x < min {
		x = min
	}
	if x > max {
		x = max
	}
	return x
}

// Remainder returns the IEEE 754 floating-point remainder of x/y.
//
// Special cases are:
//	Remainder(±Inf, y) = NaN
//	Remainder(NaN, y) = NaN
//	Remainder(x, 0) = NaN
//	Remainder(x, ±Inf) = x
//	Remainder(x, NaN) = NaN
func Remainder(a, b float32) float32 {
	return float32(math.Remainder(float64(a), float64(b)))
}

// Abs returns the absolute value of x.
//
// Special cases are:
//	Abs(±Inf) = +Inf
//	Abs(NaN) = NaN
func Abs(v float32) float32 {
	return float32(math.Abs(float64(v)))
}

// Ceil returns the least integer value greater than or equal to x.
//
// Special cases are:
//	Ceil(±0) = ±0
//	Ceil(±Inf) = ±Inf
//	Ceil(NaN) = NaN
func Ceil(v float32) float32 {
	return float32(math.Ceil(float64(v)))
}

// Floor returns the greatest integer value less than or equal to x.
//
// Special cases are:
//	Floor(±0) = ±0
//	Floor(±Inf) = ±Inf
//	Floor(NaN) = NaN
func Floor(v float32) float32 {
	return float32(math.Floor(float64(v)))
}
