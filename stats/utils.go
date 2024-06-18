package stats

import "golang.org/x/exp/constraints"

// Max compares two variables and returns the smaller one
func Max[T constraints.Ordered](a T, b T) T {
	if a < b {
		return b
	}
	return a
}
