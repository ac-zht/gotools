package list

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/zht-account/gotools"
	"testing"
)

func TestArrayList_Cap(t *testing.T) {
	testCase := []struct {
		name  string
		list  *ArrayList[int]
		index int
		want  int
	}{
		{
			name: "normal",
			list: &ArrayList[int]{
				values: make([]int, 0, 3),
			},
			want: 3,
		},
		{
			name: "nil",
			list: &ArrayList[int]{
				values: nil,
			},
			want: 0,
		},
	}
	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			c := tc.list.Cap()
			assert.Equal(t, tc.want, c)
		})
	}
}

func TestArrayList_Len(t *testing.T) {
	testCase := []struct {
		name  string
		list  *ArrayList[int]
		index int
		want  int
	}{
		{
			name: "normal",
			list: &ArrayList[int]{
				values: make([]int, 3),
			},
			want: 3,
		},
		{
			name: "nil",
			list: &ArrayList[int]{
				values: nil,
			},
			want: 0,
		},
	}
	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			l := tc.list.Len()
			assert.Equal(t, tc.want, l)
		})
	}
}

func TestArrayList_Get(t *testing.T) {
	testCase := []struct {
		name    string
		list    *ArrayList[int]
		index   int
		want    int
		wantErr error
	}{
		{
			name:  "index 1",
			list:  NewArrayListOf([]int{1, 2, 3}),
			index: 1,
			want:  2,
		},
		{
			name:    "index 3",
			list:    NewArrayListOf([]int{1, 2, 3}),
			index:   3,
			wantErr: gotools.NewErrIndexOutOfRange(3, 3),
		},
		{
			name:    "index -1",
			list:    NewArrayListOf([]int{1, 2, 3}),
			index:   -1,
			wantErr: gotools.NewErrIndexOutOfRange(3, -1),
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
		name   string
		list   *ArrayList[int]
		newVal []int
		want   []int
	}{
		//append non-empty values
		{
			name:   "append non-empty values to empty list",
			list:   NewArrayListOf[int]([]int{}),
			newVal: []int{2, 3},
			want:   []int{2, 3},
		},
		{
			name:   "append non-empty values to nil list",
			list:   NewArrayListOf[int](nil),
			newVal: []int{2, 3},
			want:   []int{2, 3},
		},
		{
			name:   "append non-empty values to non-empty list",
			list:   NewArrayListOf[int]([]int{1}),
			newVal: []int{2, 3},
			want:   []int{1, 2, 3},
		},
		//append empty values
		{
			name:   "append empty values to non-empty list",
			list:   NewArrayListOf[int]([]int{1}),
			newVal: []int{},
			want:   []int{1},
		},
		{
			name:   "append empty values to empty list",
			list:   NewArrayListOf[int]([]int{}),
			newVal: []int{},
			want:   []int{},
		},
		{
			name:   "append empty values to nil list",
			list:   NewArrayListOf[int](nil),
			newVal: []int{},
			want:   []int{},
		},
		//append nil values
		{
			name:   "append nil values to non-empty list",
			list:   NewArrayListOf[int]([]int{1}),
			newVal: nil,
			want:   []int{1},
		},
		{
			name:   "append nil values to empty list",
			list:   NewArrayListOf[int]([]int{}),
			newVal: nil,
			want:   []int{},
		},
		{
			name:   "append nil values to nil list",
			list:   NewArrayListOf[int](nil),
			newVal: nil,
			want:   []int{},
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.list.Append(tc.newVal...)
			if err != nil {
				return
			}
			assert.Equal(t, tc.want, tc.list.AsSlice())
		})
	}
}

func TestArrayList_Add(t *testing.T) {
	testCase := []struct {
		name    string
		list    *ArrayList[string]
		index   int
		val     string
		want    []string
		wantErr error
	}{
		{
			name:    "index -1",
			list:    NewArrayListOf[string]([]string{"a", "b"}),
			index:   -1,
			val:     "c",
			wantErr: gotools.NewErrIndexOutOfRange(2, -1),
		},
		{
			name:  "index 0",
			list:  NewArrayListOf[string]([]string{"a", "b"}),
			index: 0,
			val:   "c",
			want:  []string{"c", "a", "b"},
		},
		{
			name:  "index 1",
			list:  NewArrayListOf[string]([]string{"a", "b"}),
			index: 1,
			val:   "c",
			want:  []string{"a", "c", "b"},
		},
		{
			name:  "index 2",
			list:  NewArrayListOf[string]([]string{"a", "b"}),
			index: 2,
			val:   "c",
			want:  []string{"a", "b", "c"},
		},
		{
			name:    "index 3",
			list:    NewArrayListOf[string]([]string{"a", "b"}),
			index:   3,
			val:     "c",
			wantErr: gotools.NewErrIndexOutOfRange(2, 3),
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.list.Add(tc.index, tc.val)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.want, tc.list.values)
		})
	}
}

func TestArrayList_Delete(t *testing.T) {
	testCase := []struct {
		name    string
		list    *ArrayList[string]
		index   int
		want    []string
		wantVal string
		wantErr error
	}{
		{
			name:    "index out of range",
			list:    NewArrayListOf[string]([]string{"a", "b", "c"}),
			index:   3,
			wantErr: gotools.NewErrIndexOutOfRange(3, 3),
		},
		{
			name:    "deleted",
			list:    NewArrayListOf[string]([]string{"a", "b", "c"}),
			index:   1,
			want:    []string{"a", "c"},
			wantVal: "b",
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			val, err := tc.list.Delete(tc.index)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.want, tc.list.AsSlice())
			assert.Equal(t, tc.wantVal, val)
		})
	}
}

func TestArrayList_Delete_Shrink(t *testing.T) {
	testCase := []struct {
		name    string
		cap     int
		loop    int
		wantCap int
	}{
		// ----- 逻辑测试 -----
		//case 1: cap小于等于64，容量不变
		{
			name:    "case 1",
			cap:     64,
			loop:    1,
			wantCap: 64,
		},
		//case 2: cap大于64且小于等于2048，长度占容量1/4以下，缩容1/2
		{
			name:    "case 2",
			cap:     100,
			loop:    20,
			wantCap: 50,
		},
		//case 3: cap大于64且小于等于2048，长度占容量1/4以上，容量不变
		{
			name:    "case 3",
			cap:     100,
			loop:    30,
			wantCap: 100,
		},
		//case 4: cap大于2048，长度占容量1/2以下，缩容5/8，向下取整
		{
			name:    "case 4",
			cap:     2050,
			loop:    100,
			wantCap: 1281,
		},
		//case 5: cap大于2048，长度占容量1/2以上，容量不变
		{
			name:    "case 5",
			cap:     2050,
			loop:    1030,
			wantCap: 2050,
		},
		// ----- 边界测试 -----
		//case 6-1: cap65 loop 2
		{
			name:    "case 6-1",
			cap:     65,
			loop:    2,
			wantCap: 32,
		},
		//case 6-2: cap 65,loop 18 --delete--> len=17
		{
			name:    "case 6-2",
			cap:     65,
			loop:    18,
			wantCap: 65,
		},
		//case 6-3: cap 65,loop 17 --delete--> len=16
		{
			name:    "case 6-3",
			cap:     65,
			loop:    17,
			wantCap: 32,
		},
		//case 7-1: cap:2047 loop 512 --delete--> len=511
		{
			name:    "case 7-2",
			cap:     2047,
			loop:    512,
			wantCap: 1023,
		},
		//case 7-2: cap:2047 loop 513 --delete--> len=512
		{
			name:    "case 7-2",
			cap:     2047,
			loop:    513,
			wantCap: 2047,
		},
		//case 8-1: cap:2049 loop 1026 --delete--> len=1025
		{
			name:    "case 8-1",
			cap:     2049,
			loop:    1026,
			wantCap: 2049,
		},
		//case 8-2: cap:2049 loop 1025 delete --delete--> len=1024
		{
			name:    "case 8-2",
			cap:     2049,
			loop:    1025,
			wantCap: 1280,
		},
	}
	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			list := NewArrayList[int](tc.cap)
			for i := 0; i < tc.loop; i++ {
				_ = list.Append(i)
			}
			_, _ = list.Delete(0)
			assert.Equal(t, tc.wantCap, list.Cap())
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
			name:    "index -1",
			list:    NewArrayListOf([]string{"a", "b", "c"}),
			index:   -1,
			wantErr: gotools.NewErrIndexOutOfRange(3, -1),
		},
		{
			name:  "index 1",
			list:  NewArrayListOf([]string{"a", "b", "c"}),
			index: 1,
			val:   "f",
			want:  []string{"a", "f", "c"},
		},
		{
			name:    "index 3",
			list:    NewArrayListOf([]string{"a", "b", "c"}),
			index:   3,
			wantErr: gotools.NewErrIndexOutOfRange(3, 3),
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
			orgAddr := fmt.Sprintf("%p", tc.list.values)
			sliceAddr := fmt.Sprintf("%p", res)
			assert.Equal(t, tc.want, res)
			assert.NotEqual(t, orgAddr, sliceAddr)
		})
	}
}
