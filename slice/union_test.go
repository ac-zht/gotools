package slice

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUnion(t *testing.T) {
    testCase := []struct {
        name string
        s1   []string
        s2   []string
        want []string
    }{
        {
            name: "union",
            s1:   []string{"a", "b", "d"},
            s2:   []string{"b", "c"},
            want: []string{"a", "b", "d", "c"},
        },
    }

    for _, tc := range testCase {
        t.Run(tc.name, func(t *testing.T) {
            union := Union[string](tc.s1, tc.s2)
            assert.Equal(t, tc.want, union)
        })
    }
}
