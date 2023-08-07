package list

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"gotools"
	"testing"
)

func TestArrayList_Get(t *testing.T) {
	testCase := []struct {
		name    string
		list    *ArrayList[int]
		index   int
		want    int
		wantErr error
	}{
		{
			name:    "index out of range",
			list:    NewArrayListOf([]int{1, 2, 3}),
			index:   3,
			wantErr: gotools.ErrIndexOutOfRange,
		},
		{
			name:  "get",
			list:  NewArrayListOf([]int{1, 2, 3}),
			index: 1,
			want:  2,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			val, err := tc.list.Get(tc.index)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.want, val)
		})
	}
}

func TestArrayList_Append(t *testing.T) {
	testCase := []struct {
		name string
		list *ArrayList[int]
		a    int
		b    int
		want []int
	}{
		{
			name: "append",
			list: NewArrayListOf([]int{1}),
			a:    2,
			b:    3,
			want: []int{1, 2, 3},
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			_ = tc.list.Append(tc.a, tc.b)
			assert.Equal(t, tc.want, tc.list.values)
		})
	}
}

func TestArrayList_Add(t *testing.T) {
	testCase := []struct {
		name  string
		list  *ArrayList[string]
		index int
		val   string
		want  []string
	}{
		{
			name:  "add",
			list:  NewArrayListOf([]string{"a", "b", "c"}),
			index: 1,
			val:   "f",
			want:  []string{"a", "f", "b", "c"},
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			_ = tc.list.Add(tc.index, tc.val)
			assert.Equal(t, tc.want, tc.list.values)
		})
	}
}

func TestArrayList_Delete(t *testing.T) {
	testCase := []struct {
		name    string
		list    func() *ArrayList[string]
		index   int
		want    []string
		wantVal string
		wantErr error
	}{
		{
			name: "index out of range",
			list: func() *ArrayList[string] {
				return NewArrayListOf([]string{"a", "b", "c"})
			},
			index:   3,
			wantErr: gotools.ErrIndexOutOfRange,
		},
		{
			name: "no shrink",
			list: func() *ArrayList[string] {
				return NewArrayListOf([]string{"a", "b", "c"})
			},
			index:   1,
			want:    []string{"a", "c"},
			wantVal: "b",
		},
		{
			name: "shrink",
			list: func() *ArrayList[string] {
				list := NewArrayListOf(make([]string, 251, 1000))
				list.values[1] = "test"
				return list
			},
			index:   1,
			want:    make([]string, 250, 500),
			wantVal: "test",
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			list := tc.list()
			val, err := list.Delete(tc.index)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.want, list.values)
			assert.Equal(t, tc.wantVal, val)
		})
	}
}

func TestArrayList_Set(t *testing.T) {
	testCase := []struct {
		name    string
		list    *ArrayList[string]
		index   int
		val     string
		want    []string
		wantErr error
	}{
		{
			name:    "index out of range",
			list:    NewArrayListOf([]string{"a", "b", "c"}),
			index:   3,
			wantErr: gotools.ErrIndexOutOfRange,
		},
		{
			name:  "set",
			list:  NewArrayListOf([]string{"a", "b", "c"}),
			index: 1,
			val:   "f",
			want:  []string{"a", "f", "c"},
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.list.Set(tc.index, tc.val)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.want, tc.list.values)
		})
	}
}

func TestArrayList_Range(t *testing.T) {
	testCase := []struct {
		name    string
		list    *ArrayList[int]
		wantVal int
		wantErr error
	}{
		{
			name: "array is nil",
			list: &ArrayList[int]{
				values: nil,
			},
			wantVal: 0,
		},
		{
			name:    "sum",
			list:    NewArrayListOf([]int{1, 2, 3, 4}),
			wantVal: 10,
		},
		{
			name:    "test interrupt",
			list:    NewArrayListOf([]int{1, -2, 3, -4}),
			wantErr: errors.New("index 1 is error"),
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			sum := 0
			err := tc.list.Range(func(index int, num int) error {
				if num < 0 {
					return fmt.Errorf("index %d is error", index)
				}
				sum += num
				return nil
			})
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantVal, sum)
		})
	}
}

func TestArrayList_AsSlice(t *testing.T) {
	testCase := []struct {
		name string
		list *ArrayList[int]
		want []int
	}{
		{
			name: "nil",
			list: &ArrayList[int]{
				values: nil,
			},
			want: []int{},
		},
		{
			name: "int slice",
			list: &ArrayList[int]{
				values: []int{1, 2, 3},
			},
			want: []int{1, 2, 3},
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			res := tc.list.AsSlice()
			assert.Equal(t, tc.want, res)
		})
	}
}
