package diagnostics

import (
	"github.com/MontFerret/ferret/v2/pkg/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/file"
)

const (
	SyntaxError   diagnostics.Kind = "SyntaxError"
	NameError     diagnostics.Kind = "NameError"
	SemanticError diagnostics.Kind = "SemanticError"
)

func NewError(src *file.Source, kind diagnostics.Kind, message string) *diagnostics.Diagnostic {
	return &diagnostics.Diagnostic{
		Message: message,
		Source:  src,
		Kind:    kind,
	}
}

func NewUnexpectedError(src *file.Source, message string) *diagnostics.Diagnostic {
	return NewError(src, diagnostics.UnexpectedError, message)
}

func NewUnexpectedErrorWith(src *file.Source, message string, cause error) *diagnostics.Diagnostic {
	e := NewUnexpectedError(src, message)
	e.Cause = cause

	return e
}

func NewEmptyQueryError(src *file.Source) *diagnostics.Diagnostic {
	return NewError(src, SyntaxError, "Query is empty")
}
