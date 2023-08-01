package slice

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWeightRandomIndex(t *testing.T) {
	testCase := []struct {
		name string
		src  []int
		want int
	}{
		{
			name: "index 0",
			src:  []int{1, 0},
			want: 0,
		},
		{
			name: "index 1",
			src:  []int{0, 1},
			want: 1,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			index := WeightRandomIndex(tc.src)
			assert.Equal(t, tc.want, index)
		})
	}
	t.Run("random index", func(t *testing.T) {
		src := []int{1, 1}
		index := WeightRandomIndex(src)
		assert.GreaterOrEqual(t, index, 0)
		assert.Less(t, index, len(src))
	})
}
