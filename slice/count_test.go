package slice

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCountValues_IntSlice(t *testing.T) {
	testCase := []struct {
		name string
		src  []int
		want map[int]int
	}{
		{
			name: "int slice",
			src:  []int{1, 2, 2, 1, 1, 3},
			want: map[int]int{
				1: 3,
				2: 2,
				3: 1,
			},
		},
	}
	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			mp := CountValues[int](tc.src)
			assert.Equal(t, tc.want, mp)
		})
	}
}
