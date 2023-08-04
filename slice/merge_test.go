package slice

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMerge(t *testing.T) {
	testCase := []struct {
		name string
		s1   []int
		s2   []int
		s3   []int
		want []int
	}{
		{
			name: "nil",
			s1:   nil,
			s2:   nil,
			s3:   nil,
			want: []int{},
		},
		{
			name: "merge",
			s1:   []int{1, 2},
			s2:   []int{3, 4},
			s3:   []int{5, 6},
			want: []int{1, 2, 3, 4, 5, 6},
		},
	}
	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			ret := Merge[int](tc.s1, tc.s2, tc.s3)
			assert.Equal(t, tc.want, ret)
		})
	}
}
