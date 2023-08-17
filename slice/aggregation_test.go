package slice

import (
	"github.com/stretchr/testify/assert"
	"github.com/zht-account/gotools"
	"testing"
)

func TestMax(t *testing.T) {
	t.Run("panic input", func(t *testing.T) {
		assert.Panics(t, func() {
			Max[int](nil)
		})
		assert.Panics(t, func() {
			Max[int]([]int{})
		})
	})

	testCase := []struct {
		name  string
		input []int
		want  int
	}{
		{
			name:  "value",
			input: []int{1},
			want:  1,
		},
		{
			name:  "values",
			input: []int{1, 2, 3},
			want:  3,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			res := Max[int](tc.input)
			assert.Equal(t, tc.want, res)
		})
	}

	testMaxTypes[uint](t)
	testMaxTypes[uint8](t)
	testMaxTypes[uint16](t)
	testMaxTypes[uint32](t)
	testMaxTypes[uint64](t)
	testMaxTypes[int8](t)
	testMaxTypes[int16](t)
	testMaxTypes[int32](t)
	testMaxTypes[int64](t)
	testMaxTypes[float32](t)
	testMaxTypes[float64](t)
}

func testMaxTypes[T gotools.RealNumber](t *testing.T) {
	t.Run("testMaxTypes", func(t *testing.T) {
		res := Max[T]([]T{1, 2, 3})
		assert.Equal(t, T(3), res)
	})
}

func TestMin(t *testing.T) {
	t.Run("panic input", func(t *testing.T) {
		assert.Panics(t, func() {
			Min[int](nil)
		})
		assert.Panics(t, func() {
			Min[int]([]int{})
		})
	})

	testCase := []struct {
		name  string
		input []int
		want  int
	}{
		{
			name:  "value",
			input: []int{1},
			want:  1,
		},
		{
			name:  "values",
			input: []int{1, 2, 3},
			want:  1,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			res := Min[int](tc.input)
			assert.Equal(t, tc.want, res)
		})
	}

	testMinTypes[uint](t)
	testMinTypes[uint8](t)
	testMinTypes[uint16](t)
	testMinTypes[uint32](t)
	testMinTypes[uint64](t)
	testMinTypes[int8](t)
	testMinTypes[int16](t)
	testMinTypes[int32](t)
	testMinTypes[int64](t)
	testMinTypes[float32](t)
	testMinTypes[float64](t)
}

func testMinTypes[T gotools.RealNumber](t *testing.T) {
	t.Run("testMinTypes", func(t *testing.T) {
		res := Min[T]([]T{1, 2, 3})
		assert.Equal(t, T(1), res)
	})
}

func TestSum(t *testing.T) {
	testCase := []struct {
		name  string
		input []int
		want  int
	}{
		{
			name: "nil",
		},
		{
			name:  "empty",
			input: []int{},
		},
		{
			name:  "value",
			input: []int{1},
			want:  1,
		},
		{
			name:  "values",
			input: []int{1, 2, 3},
			want:  6,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			res := Sum[int](tc.input)
			assert.Equal(t, tc.want, res)
		})
	}

	testSumTypes[uint](t)
	testSumTypes[uint8](t)
	testSumTypes[uint16](t)
	testSumTypes[uint32](t)
	testSumTypes[uint64](t)
	testSumTypes[int8](t)
	testSumTypes[int16](t)
	testSumTypes[int32](t)
	testSumTypes[int64](t)
	testSumTypes[float32](t)
	testSumTypes[float64](t)
	testSumTypes[complex64](t)
	testSumTypes[complex128](t)
}

func testSumTypes[T gotools.Number](t *testing.T) {
	t.Run("testSumTypes", func(t *testing.T) {
		res := Sum[T]([]T{1, 2, 3})
		assert.Equal(t, T(6), res)
	})
}
