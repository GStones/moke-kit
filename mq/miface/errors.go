package miface

import "errors"

var (
	// ErrDataAlreadySet is returned when a publish payload is set more than once.
	ErrDataAlreadySet = errors.New("ErrDataAlreadySet")
	// ErrSemanticsAlreadySet is returned when delivery semantics are set more than once.
	ErrSemanticsAlreadySet = errors.New("ErrSemanticsAlreadySet")
)
