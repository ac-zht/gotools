package slice

import (
	"github.com/stretchr/testify/assert"
	"gotools"
	"testing"
)

func TestDelete(t *testing.T) {
	testCase := []struct {
		name      string
		index     int
		s         []int
		wantSlice []int
		wantVal   int
		wantErr   error
	}{
		{
			name:      "delete",
			index:     2,
			s:         []int{1, 2, 3, 4},
			wantSlice: []int{1, 2, 4},
			wantVal:   3,
		},
		{
			name:    "index out of range",
			index:   4,
			s:       []int{1, 2, 3, 4},
			wantErr: gotools.ErrIndexOutOfRange,
		},
	}
	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			res, val, err := Delete[int](tc.index, tc.s)
			if err != nil {
				assert.Equal(t, tc.wantErr, err)
				return
			}
			assert.Equal(t, tc.wantVal, val)
			assert.Equal(t, tc.wantSlice, res)
		})
	}
}

func TestFilterDelete(t *testing.T) {
	testCase := []struct {
		name   string
		s      []int
		filter func(key int, value int) bool
		want   []int
	}{
		{
			name: "key filter delete",
			s:    []int{1, 2, 3, 2},
			filter: func(key, value int) bool {
				if key == 1 {
					return true
				}
				return false
			},
			want: []int{1, 3, 2},
		},
		{
			name: "value filter delete",
			s:    []int{1, 2, 3, 2},
			filter: func(key, value int) bool {
				if value == 2 {
					return true
				}
				return false
			},
			want: []int{1, 3},
		},
	}
	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			s := FilterDelete[int](tc.s, tc.filter)
			assert.Equal(t, tc.want, s)
		})
	}
}
