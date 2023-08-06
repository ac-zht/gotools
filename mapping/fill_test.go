package mapping

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFillKeys_KeyStringValueInt(t *testing.T) {
	testCase := []struct {
		name string
		keys []string
		val  int
		want map[string]int
	}{
		{
			name: "fill key string value int",
			keys: []string{"a", "b", "c"},
			val:  1,
			want: map[string]int{
				"a": 1,
				"b": 1,
				"c": 1,
			},
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			fill := FillKeys(tc.keys, tc.val)
			assert.Equal(t, tc.want, fill)
		})
	}

	TestFillKeys_KeyStringValueMap(t)
}

func TestFillKeys_KeyStringValueMap(t *testing.T) {
	t.Run("fill key string value map", func(t *testing.T) {
		fill := FillKeys[string, map[string]string]([]string{"first", "second", "something"}, map[string]string{
			"a": "first",
			"b": "second",
			"c": "something",
		})
		assert.Equal(t, map[string]map[string]string{
			"first": {
				"a": "first",
				"b": "second",
				"c": "something",
			},
			"second": {
				"a": "first",
				"b": "second",
				"c": "something",
			},
			"something": {
				"a": "first",
				"b": "second",
				"c": "something",
			},
		}, fill)
	})
}
