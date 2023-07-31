package slice

import (
    "github.com/stretchr/testify/assert"
    "testing"
)

func TestDiff(t *testing.T) {
    testCase := []struct {
        name string
        src  []int
        dst  []int
        want []int
    }{
        {
            name: "exist diff",
            src:  []int{1, 2},
            dst:  []int{2, 3, 4},
            want: []int{2},
        },

        {
            name: "no diff",
            src:  []int{1, 2},
            dst:  []int{3, 4},
            want: []int{},
        },
    }

    for _, tc := range testCase {
        t.Run(tc.name, func(t *testing.T) {
            d := Diff[int](tc.src, tc.dst)
            assert.Equal(t, tc.want, d)
        })
    }
}
