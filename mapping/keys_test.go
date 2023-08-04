package mapping

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestKeys(t *testing.T) {
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
			want: []string{"a", "b", "c", "d"},
		},
	}
	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			keys := Keys[string, string](tc.mp)
			assert.ElementsMatch(t, tc.want, keys)
		})
	}
}

func TestKeyExist(t *testing.T) {
	testCase := []struct {
		name string
		mp   map[string]string
		key  string
		want bool
	}{
		{
			name: "key exist",
			mp: map[string]string{
				"a": "w",
				"b": "x",
			},
			key:  "a",
			want: true,
		},
		{
			name: "key not exist",
			mp: map[string]string{
				"a": "w",
				"b": "x",
			},
			key:  "c",
			want: false,
		},
	}
	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			exist := KeyExist[string](tc.mp, tc.key)
			assert.Equal(t, tc.want, exist)
		})
	}
}

func TestKeysWithValue(t *testing.T) {
	testCase := []struct {
		name string
		mp   map[string]int
		val  int
		want []string
	}{
		{
			name: "keys with value",
			mp: map[string]int{
				"a": 1,
				"b": 2,
				"c": 1,
			},
			val:  1,
			want: []string{"a", "c"},
		},
	}
	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			keys := KeysWithValue[string, int](tc.mp, tc.val)
			assert.ElementsMatch(t, tc.want, keys)
		})
	}
}
