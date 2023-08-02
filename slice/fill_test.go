package slice

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFill(t *testing.T) {
	testCase := []struct {
		name  string
		start int
		cnt   int
		val   string
		want  []string
	}{
		{
			name:  "start = 0 string slice",
			start: 0,
			cnt:   2,
			val:   "a",
			want:  []string{"a", "a"},
		},
		{
			name:  "start > 0 string slice",
			start: 1,
			cnt:   2,
			val:   "a",
			want:  []string{"", "a", "a"},
		},
	}
	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			fill := Fill[string](tc.start, tc.cnt, tc.val)
			assert.Equal(t, tc.want, fill)
		})
	}
}
