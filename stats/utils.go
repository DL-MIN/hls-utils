package stats

import "golang.org/x/exp/constraints"

func Min[T constraints.Ordered](a T, b T) T {
    if a > b {
        return b
    }
    return a
}
