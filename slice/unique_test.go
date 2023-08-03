package slice

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUnique(t *testing.T) {
	testCase := []struct {
		name string
		src  []int
		want []int
	}{
		{
			name: "empty",
			src:  []int{},
			want: []int{},
		},
		{
			name: "unique",
			src:  []int{1, 2, 2, 1},
			want: []int{1, 2},
		},
	}
	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			uniq := Unique[int](tc.src)
			assert.ElementsMatch(t, tc.want, uniq)
		})
	}
}
