package gotools

import (
	"errors"
	"fmt"
)

var (
	ErrSliceLengthNotEnough = errors.New("gotools: slice not enough length")
	ErrSliceIsEmpty         = errors.New("gotools: slice length is empty")
	ErrIndexOutOfRange      = errors.New("gotools: index out of slice range")
)

func NewErrIndexOutOfRange(length int, index int) error {
	return fmt.Errorf("gotools: index out of range, len %d, index %d", length, index)
}
