package list

import (
	"gotools"
	"gotools/slice"
)

type ArrayList[T any] struct {
	values []T
}

func NewArrayList[T any](cap int) *ArrayList[T] {
	return &ArrayList[T]{
		values: make([]T, 0, cap),
	}
}

func NewArrayListOf[T any](values []T) *ArrayList[T] {
	return &ArrayList[T]{
		values: values,
	}
}

func (a *ArrayList[T]) Len() int {
	return len(a.values)
}

func (a *ArrayList[T]) Cap() int {
	return cap(a.values)
}

func (a *ArrayList[T]) Get(index int) (t T, err error) {
	if index < 0 || index >= len(a.values) {
		return t, gotools.NewErrIndexOutOfRange(len(a.values), index)
	}
	return a.values[index], nil
}

func (a *ArrayList[T]) Append(t ...T) error {
	a.values = append(a.values, t...)
	return nil
}

func (a *ArrayList[T]) Add(index int, t T) error {
	if index < 0 || index > len(a.values) {
		return gotools.NewErrIndexOutOfRange(len(a.values), index)
	}
	a.values = append(a.values, t)
	copy(a.values[index+1:], a.values[index:])
	a.values[index] = t
	return nil
}

func (a *ArrayList[T]) Delete(index int) (t T, err error) {
	a.values, t, err = slice.Delete(index, a.values)
	if err != nil {
		return t, err
	}
	a.Shrink()
	return t, nil
}

func (a *ArrayList[T]) Shrink() {
	a.values = slice.Shrink[T](a.values)
}

func (a *ArrayList[T]) Set(index int, t T) error {
	if index < 0 || index >= len(a.values) {
		return gotools.NewErrIndexOutOfRange(len(a.values), index)
	}
	a.values[index] = t
	return nil
}

func (a *ArrayList[T]) Range(fn func(index int, t T) error) error {
	for idx, val := range a.values {
		err := fn(idx, val)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *ArrayList[T]) AsSlice() []T {
	res := make([]T, len(a.values))
	copy(res, a.values)
	return res
}
