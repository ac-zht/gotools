package slice

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestShuffleString(t *testing.T) {
	testCase := []struct {
		name string
		src  []string
		want []string
	}{
		{
			name: "nil",
			src:  nil,
			want: nil,
		},
		{
			name: "empty",
			src:  []string{},
			want: []string{},
		},
		{
			name: "shuffle",
			src:  []string{"a", "b", "c", "d"},
			want: []string{"a", "b", "c", "d"},
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			Shuffle[string](tc.src)
			assert.ElementsMatch(t, tc.want, tc.src)
		})
	}
}
