package slice

import (
	"github.com/stretchr/testify/assert"
	"gotools"
	"testing"
)

func TestRandomIndexes(t *testing.T) {
	testCase := []struct {
		name    string
		src     []int
		n       int
		want    []int
		wantErr error
	}{
		{
			name:    "slice empty",
			src:     []int{},
			n:       3,
			wantErr: gotools.ErrSliceIsEmpty,
		},
		{
			name:    "n > len(src)",
			src:     []int{1, 2},
			n:       3,
			wantErr: gotools.ErrSliceLengthNotEnough,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			index, err := RandomIndexes(tc.src, tc.n)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.want, index)
		})
	}
	t.Run("random index", func(t *testing.T) {
		src := []string{"a", "b", "c", "d"}
		indexes, _ := RandomIndexes[string](src, 2)
		for _, v := range indexes {
			assert.GreaterOrEqual(t, v, 0)
			assert.Less(t, v, len(src))
		}
	})
}

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
