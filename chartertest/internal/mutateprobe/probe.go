// Package mutateprobe is code for the mutate tests to break: a function
// whose test notices, one whose test does not, and a variable no mutant
// can change. The gate ignores it; the tests describe its anchors to
// mutate themselves.
package mutateprobe

// Clamp holds v between lo and hi.
//
//charter:realizes SW-001
func Clamp(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

// Label names a count.
//
//charter:realizes SW-001
func Label(n int) string {
	if n == 1 {
		return "one"
	}
	return "many"
}

// Limit is a number no mutant changes.
//
//charter:realizes SW-001
var Limit = 3
