package slice

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDiff(t *testing.T) {
	testCase := []struct {
		name string
		src  []int
		dst  []int
		want []int
	}{
		{
			name: "diff",
			src:  []int{1, 2, 3},
			dst:  []int{2, 3, 4},
			want: []int{1},
		},
		{
			name: "no diff",
			src:  []int{1, 2, 3},
			dst:  []int{1, 2, 3},
			want: []int{},
		},
	}
	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			diff := Diff[int](tc.src, tc.dst)
			assert.Equal(t, tc.want, diff)
		})
	}
}

func TestSymmetricDiff(t *testing.T) {
	testCase := []struct {
		name string
		src  []int
		dst  []int
		want []int
	}{
		{
			name: "symmetric diff",
			src:  []int{1, 2, 3},
			dst:  []int{2, 3, 4},
			want: []int{1, 4},
		},
		{
			name: "no diff",
			src:  []int{1, 2, 3},
			dst:  []int{1, 2, 3},
			want: []int{},
		},
	}
	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			diff := SymmetricDiff[int](tc.src, tc.dst)
			assert.ElementsMatch(t, tc.want, diff)
		})
	}
}
