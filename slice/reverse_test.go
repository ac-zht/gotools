package slice

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReverseInt(t *testing.T) {
	testCase := []struct {
		name string
		src  []int
		want []int
	}{
		{
			name: "nil",
			src:  nil,
			want: []int{},
		},
		{
			name: "empty",
			src:  []int{},
			want: []int{},
		},
		{
			name: "reverse",
			src:  []int{1, 2, 3},
			want: []int{3, 2, 1},
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			reverse := Reverse[int](tc.src)
			assert.Equal(t, tc.want, reverse)
		})
	}
}

func TestReverseSelfInt(t *testing.T) {
	testCase := []struct {
		name string
		src  []int
		want []int
	}{
		{
			name: "nil",
			src:  nil,
			want: nil,
		},
		{
			name: "empty",
			src:  []int{},
			want: []int{},
		},
		{
			name: "reverse",
			src:  []int{1, 2, 3},
			want: []int{3, 2, 1},
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			ReverseSelf[int](tc.src)
			assert.Equal(t, tc.want, tc.src)
		})
	}
}
