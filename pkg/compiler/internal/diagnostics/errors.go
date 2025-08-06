package diagnostics

import (
	"github.com/MontFerret/ferret/pkg/file"
)

var (
	ErrNotImplemented    = "not implemented"
	ErrInvalidToken      = "invalid token"
	ErrConstantNotFound  = "constant not found"
	ErrInvalidDataSource = "invalid data source"
	ErrUnknownOpcode     = "unknown opcode"
)

func NewEmptyQueryErr(src *file.Source) *CompilationError {
	return &CompilationError{
		Message: "Query is empty",
		Source:  src,
		Kind:    SyntaxError,
	}
}

func NewInternalErr(src *file.Source, msg string) *CompilationError {
	return &CompilationError{
		Message: msg,
		Source:  src,
		Kind:    InternalError,
	}
}

func NewInternalErrWith(src *file.Source, msg string, cause error) *CompilationError {
	return &CompilationError{
		Message: msg,
		Source:  src,
		Kind:    InternalError,
		Cause:   cause,
	}
}
