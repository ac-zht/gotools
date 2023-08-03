package slice

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIntersect(t *testing.T) {
	testCase := []struct {
		name string
		src  []int
		dst  []int
		want []int
	}{
		{
			name: "exist intersection",
			src:  []int{1, 2},
			dst:  []int{2, 3, 4},
			want: []int{2},
		},
		{
			name: "contain",
			src:  []int{1, 2},
			dst:  []int{1, 2, 3},
			want: []int{1, 2},
		},
		{
			name: "no intersection",
			src:  []int{1, 2},
			dst:  []int{3, 4},
			want: []int{},
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			intersect := Intersect[int](tc.src, tc.dst)
			assert.Equal(t, tc.want, intersect)
		})
	}
}
