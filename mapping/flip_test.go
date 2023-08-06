package mapping

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFlip(t *testing.T) {
	testCase := []struct {
		name string
		src  map[string]int
		want map[int]string
	}{
		{
			name: "unique value map flip",
			src: map[string]int{
				"a": 1,
				"b": 2,
				"c": 3,
			},
			want: map[int]string{
				1: "a",
				2: "b",
				3: "c",
			},
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			flip := Flip(tc.src)
			assert.Equal(t, tc.want, flip)
		})
	}
}
