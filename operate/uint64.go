package operate

import "math"

// Uint64Add
func Uint64Add(x, y uint64) uint64 {
	if addOverflow(x, y) {
		return math.MaxUint64
	}
	return x + y
}

// Uint64Sub
func Uint64Sub(x, y uint64) uint64 {
	if subOverflow(x, y) {
		return 0
	}
	return x - y
}

// Uint64Inc
func Uint64Inc(x uint64) uint64 {
	return Uint64Add(x, 1)
}

// Uint64Dec
func Uint64Dec(x uint64) uint64 {
	return Uint64Sub(x, 1)
}

func addOverflow(x, y uint64) bool {
	return x+y < x
}

func subOverflow(x, y uint64) bool {
	return x < y
}
