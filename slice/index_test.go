package slice

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIndex(t *testing.T) {
	testCase := []struct {
		name string
		val  int
		s    []int
		want int
	}{
		{
			name: "find",
			val:  2,
			s:    []int{1, 2, 3, 2},
			want: 1,
		},
		{
			name: "no find",
			val:  4,
			s:    []int{1, 2, 3},
			want: -1,
		},
	}
	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			index := Index[int](tc.val, tc.s)
			assert.Equal(t, tc.want, index)
		})
	}
}

func TestLastIndex(t *testing.T) {
	testCase := []struct {
		name string
		val  int
		s    []int
		want int
	}{
		{
			name: "find",
			val:  2,
			s:    []int{1, 2, 3, 2},
			want: 3,
		},
		{
			name: "no find",
			val:  4,
			s:    []int{1, 2, 3},
			want: -1,
		},
	}
	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			index := LastIndex[int](tc.val, tc.s)
			assert.Equal(t, tc.want, index)
		})
	}
}

func TestAllIndex(t *testing.T) {
	testCase := []struct {
		name string
		val  int
		s    []int
		want []int
	}{
		{
			name: "find",
			val:  2,
			s:    []int{1, 2, 3, 2},
			want: []int{1, 3},
		},
		{
			name: "no find",
			val:  4,
			s:    []int{1, 2, 3},
			want: []int{},
		},
	}
	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			indexes := AllIndex[int](tc.val, tc.s)
			assert.Equal(t, tc.want, indexes)
		})
	}
}
