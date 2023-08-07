package slice

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestShrink(t *testing.T) {
	testCase := []struct {
		name    string
		src     []any
		want    []any
		wantCap int
	}{
		{
			name:    "c <= 64",
			src:     make([]any, 45, 60),
			want:    make([]any, 45, 60),
			wantCap: 60,
		},
		{
			name:    "c > 64 and c <= 2048 and l <= c/4",
			src:     make([]any, 220, 1000),
			want:    make([]any, 220, 1000),
			wantCap: 500,
		},
		{
			name:    "c > 2048 and l < c/2",
			src:     make([]any, 1000, 3000),
			want:    make([]any, 1000, 3000),
			wantCap: 1875,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			res := Shrink(tc.src)
			assert.Equal(t, tc.want, res)
			assert.Equal(t, tc.wantCap, cap(res))
		})
	}
}
