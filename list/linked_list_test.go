package list

import (
    "errors"
    "fmt"
    "github.com/stretchr/testify/assert"
    "gotools"
    "testing"
)

func TestLinkedList_Get(t *testing.T) {
    testCase := []struct {
        name      string
        list      *LinkedList[int]
        index     int
        want      int
        wantError error
    }{
        {
            name:      "index -1",
            list:      NewLinkedListOf[int]([]int{1, 2, 3, 4}),
            index:     -1,
            wantError: gotools.NewErrIndexOutOfRange(4, -1),
        },
        {
            name:  "index 2",
            list:  NewLinkedListOf[int]([]int{1, 2, 3, 4}),
            index: 2,
            want:  3,
        },
        {
            name:      "index 4",
            list:      NewLinkedListOf[int]([]int{1, 2, 3, 4}),
            index:     4,
            wantError: gotools.NewErrIndexOutOfRange(4, 4),
        },
    }
    for _, tc := range testCase {
        t.Run(tc.name, func(t *testing.T) {
            val, err := tc.list.Get(tc.index)
            assert.Equal(t, tc.wantError, err)
            if err != nil {
                return
            }
            assert.Equal(t, tc.want, val)
        })
    }
}

func TestLinkedList_Add(t *testing.T) {
    testCase := []struct {
        name      string
        list      *LinkedList[int]
        index     int
        val       int
        wantSlice []int
        wantError error
    }{
        {
            name:      "index -1",
            list:      NewLinkedListOf[int]([]int{1, 2, 3}),
            index:     -1,
            wantError: gotools.NewErrIndexOutOfRange(3, -1),
        },
        {
            name:      "index 1",
            list:      NewLinkedListOf[int]([]int{1, 2, 3}),
            index:     1,
            val:       4,
            wantSlice: []int{1, 2, 4, 3},
        },
        {
            name:      "index 3",
            list:      NewLinkedListOf[int]([]int{1, 2, 3}),
            index:     3,
            wantError: gotools.NewErrIndexOutOfRange(3, 3),
        },
    }
    for _, tc := range testCase {
        t.Run(tc.name, func(t *testing.T) {
            err := tc.list.Add(tc.index, tc.val)
            assert.Equal(t, tc.wantError, err)
            if err != nil {
                return
            }
            assert.Equal(t, tc.wantSlice, tc.list.AsSlice())
            assert.Equal(t, len(tc.wantSlice), tc.list.Len())
        })
    }
}

func TestLinkedList_Append(t *testing.T) {
    testCase := []struct {
        name      string
        list      *LinkedList[int]
        input     []int
        wantSlice []int
        wantError error
    }{
        {
            name:      "normal",
            list:      NewLinkedListOf[int]([]int{1, 2, 3}),
            input:     []int{4, 5},
            wantSlice: []int{1, 2, 3, 4, 5},
        },
        {
            name:      "nil",
            list:      NewLinkedListOf[int]([]int{1, 2, 3}),
            input:     nil,
            wantSlice: []int{1, 2, 3},
        },
    }
    for _, tc := range testCase {
        t.Run(tc.name, func(t *testing.T) {
            tc.list.Append(tc.input...)
            assert.Equal(t, tc.wantSlice, tc.list.AsSlice())
            assert.Equal(t, len(tc.wantSlice), tc.list.Len())
        })
    }
}

func TestLinkedList_Delete(t *testing.T) {
    testCase := []struct {
        name      string
        list      *LinkedList[int]
        index     int
        wantVal   int
        wantSlice []int
        wantError error
    }{
        {
            name:      "index 2 deleted",
            list:      NewLinkedListOf[int]([]int{1, 2, 3, 4, 5}),
            index:     2,
            wantVal:   3,
            wantSlice: []int{1, 2, 4, 5},
        },
        {
            name:      "index 5",
            list:      NewLinkedListOf[int]([]int{1, 2, 3, 4, 5}),
            index:     5,
            wantError: gotools.NewErrIndexOutOfRange(5, 5),
        },
    }
    for _, tc := range testCase {
        t.Run(tc.name, func(t *testing.T) {
            val, err := tc.list.Delete(tc.index)
            assert.Equal(t, tc.wantError, err)
            if err != nil {
                return
            }
            assert.Equal(t, tc.wantVal, val)
            assert.Equal(t, tc.wantSlice, tc.list.AsSlice())
            assert.Equal(t, len(tc.wantSlice), tc.list.Len())
        })
    }
}

func TestLinkedList_Set(t *testing.T) {
    testCase := []struct {
        name      string
        list      *LinkedList[int]
        index     int
        value     int
        wantSlice []int
        wantError error
    }{
        {
            name:      "index 2",
            list:      NewLinkedListOf[int]([]int{1, 2, 3, 4, 5}),
            index:     2,
            value:     100,
            wantSlice: []int{1, 2, 100, 4, 5},
        },
        {
            name:      "index 5",
            list:      NewLinkedListOf[int]([]int{1, 2, 3, 4, 5}),
            index:     5,
            wantError: gotools.NewErrIndexOutOfRange(5, 5),
        },
    }
    for _, tc := range testCase {
        t.Run(tc.name, func(t *testing.T) {
            err := tc.list.Set(tc.index, tc.value)
            assert.Equal(t, tc.wantError, err)
            if err != nil {
                return
            }
            assert.Equal(t, tc.wantSlice, tc.list.AsSlice())
        })
    }
}

func TestLinkedList_Range(t *testing.T) {
    testCase := []struct {
        name      string
        list      *LinkedList[int]
        want      int
        wantError error
    }{
        {
            name: "nil",
            list: NewLinkedListOf[int](nil),
            want: 0,
        },
        {
            name: "empty",
            list: NewLinkedListOf[int]([]int{}),
            want: 0,
        },
        {
            name:      "error interrupt",
            list:      NewLinkedListOf[int]([]int{1, 2, -2, 2, 1}),
            wantError: errors.New("index 2 is error"),
        },
        {
            name: "sum",
            list: NewLinkedListOf[int]([]int{1, 2, 3, 4, 5}),
            want: 15,
        },
    }
    for _, tc := range testCase {
        t.Run(tc.name, func(t *testing.T) {
            sum := 0
            err := tc.list.Range(func(i int, val int) error {
                if val < 0 {
                    return fmt.Errorf("index %d is error", i)
                }
                sum += val
                return nil
            })
            assert.Equal(t, tc.wantError, err)
            if err != nil {
                return
            }
            assert.Equal(t, tc.want, sum)
        })
    }
}

func TestLinkedList_AsSlice(t *testing.T) {
    list := NewLinkedListOf[int]([]int{1, 2, 3, 4})
    t.Run("as slice", func(t *testing.T) {
        slice := list.AsSlice()
        assert.Equal(t, []int{1, 2, 3, 4}, slice)
    })
}
