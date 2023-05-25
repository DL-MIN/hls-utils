package stats

import "golang.org/x/exp/constraints"

// Min compares two variables and returns the smaller one
func Min[T constraints.Ordered](a T, b T) T {
	if a > b {
		return b
	}
	return a
}
