package slice

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCut(t *testing.T) {
	testCase := []struct {
		name   string
		src    []string
		length int
		want   [][]string
	}{
		{
			name:   "slice length less than cut length",
			src:    []string{"a"},
			length: 3,
			want: [][]string{
				{
					"a",
				},
			},
		},
		{
			name:   "slice length equal to cut length",
			src:    []string{"a"},
			length: 1,
			want: [][]string{
				{
					"a",
				},
			},
		},
		{
			name:   "average cut",
			src:    []string{"a", "b", "c"},
			length: 1,
			want: [][]string{
				{
					"a",
				},
				{
					"b",
				},
				{
					"c",
				},
			},
		},
		{
			name:   "uneven cut",
			src:    []string{"a", "b", "c", "d", "e"},
			length: 2,
			want: [][]string{
				{
					"a",
					"b",
				},
				{
					"c",
					"d",
				},
				{
					"e",
				},
			},
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			ret := Cut[string](tc.src, tc.length)
			assert.Equal(t, tc.want, ret)
		})
	}
}
