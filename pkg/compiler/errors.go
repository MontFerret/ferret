package compiler

import "github.com/pkg/errors"

var (
	ErrEmptyQuery        = errors.New("empty query")
	ErrNotImplemented    = errors.New("not implemented")
	ErrVariableNotFound  = errors.New("variable not found")
	ErrVariableNotUnique = errors.New("variable is already defined")
	ErrInvalidToken      = errors.New("invalid token")
	ErrUnexpectedToken   = errors.New("unexpected token")
	ErrInvalidDataSource = errors.New("invalid data source")
)
