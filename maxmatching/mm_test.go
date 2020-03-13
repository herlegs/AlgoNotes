package maxmatching

import "testing"

func TestMaxMatching(t *testing.T) {
	MaxMatching([][]int{
		{0, 1, 2, 3},
		{2, 3},
		{0},
		{2},
	}, 4)
}
