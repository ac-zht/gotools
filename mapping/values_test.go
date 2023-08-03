package mapping

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValues(t *testing.T) {
    testCase := []struct {
        name string
        mp   map[string]string
        want []string
    }{
        {
            name: "string map",
            mp: map[string]string{
                "a": "w",
                "b": "x",
                "c": "y",
                "d": "z",
            },
            want: []string{"w", "x", "y", "z"},
        },
    }
    for _, tc := range testCase {
        t.Run(tc.name, func(t *testing.T) {
            values := Values[string, string](tc.mp)
            assert.ElementsMatch(t, tc.want, values)
        })
    }
}
