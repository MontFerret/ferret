package vm

import "errors"

var (
	ErrMissedParam      = errors.New("missed parameter")
	ErrFunctionNotFound = errors.New("function not found")
)
