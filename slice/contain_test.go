package slice

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestContain(t *testing.T) {
	testCase := []struct {
		name string
		s    []int
		val  int
		want bool
	}{
		{
			name: "contain",
			s:    []int{1, 2, 3},
			val:  1,
			want: true,
		},
		{
			name: "no contain",
			s:    []int{1, 2, 3},
			val:  4,
			want: false,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			contain := Contain[int](tc.s, tc.val)
			assert.Equal(t, tc.want, contain)
		})
	}
}
