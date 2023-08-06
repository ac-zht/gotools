package slice

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCombine(t *testing.T) {
	testCase := []struct {
		name   string
		keys   []string
		values []int
		want   map[string]int
	}{
		{
			name:   "combine",
			keys:   []string{"a", "b", "c"},
			values: []int{1, 2, 3},
			want: map[string]int{
				"a": 1,
				"b": 2,
				"c": 3,
			},
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			combine := Combine(tc.keys, tc.values)
			assert.Equal(t, tc.want, combine)
		})
	}
}
