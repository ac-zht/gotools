package slice

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMap(t *testing.T) {
	testCase := []struct {
		name  string
		input []int
		want  map[int]struct{}
	}{
		{
			name:  "int slice",
			input: []int{1, 2, 3},
			want: map[int]struct{}{
				1: {},
				2: {},
				3: {},
			},
		},
	}
	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			m := Map[int](tc.input)
			assert.Equal(t, tc.want, m)
		})
	}
}
