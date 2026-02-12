package diagnostics

import (
	"github.com/MontFerret/ferret/v2/pkg/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/file"
)

const (
	ErrNotImplemented    = "not implemented"
	ErrInvalidToken      = "invalid token"
	ErrConstantNotFound  = "constant not found"
	ErrInvalidDataSource = "invalid data source"
	ErrUnknownOpcode     = "unknown opcode"
)

const (
	SyntaxError   diagnostics.Kind = "SyntaxError"
	NameError     diagnostics.Kind = "NameError"
	SemanticError diagnostics.Kind = "SemanticError"
)

type CompilationError struct {
	*diagnostics.Diagnostic
}

func NewError(src *file.Source, kind diagnostics.Kind, message string) *CompilationError {
	return &CompilationError{
		Diagnostic: &diagnostics.Diagnostic{
			Message: message,
			Source:  src,
			Kind:    kind,
		},
	}
}

func NewUnexpectedError(src *file.Source, message string) *CompilationError {
	return NewError(src, diagnostics.UnexpectedError, message)
}

func NewUnexpectedErrorWith(src *file.Source, message string, cause error) *CompilationError {
	e := NewUnexpectedError(src, message)
	e.Cause = cause

	return e
}

func NewEmptyQueryError(src *file.Source) *CompilationError {
	return NewError(src, SyntaxError, "Query is empty")
}
