package gotools

import "errors"

var (
	ErrSliceLengthNotEnough = errors.New("gotools: slice not enough length")
	ErrSliceIsEmpty         = errors.New("gotools: slice length is empty")
	ErrIndexOutOfRange      = errors.New("gotools: index out of slice range")
)
