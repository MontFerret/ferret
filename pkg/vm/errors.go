package vm

import "errors"

var (
	ErrMissedParam           = errors.New("missed parameter")
	ErrFunctionNotFound      = errors.New("function not found")
	ErrInsufficientRegisters = errors.New("insufficient registers")
	ErrUnresolvedFunction    = errors.New("unresolved function")
	ErrInvalidFunctionName   = errors.New("invalid function name")
)
