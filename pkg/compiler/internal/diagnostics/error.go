package diagnostics

import (
	"strings"

	"github.com/MontFerret/ferret/pkg/file"
)

type (
	ErrorKind string

	CompilationError struct {
		Kind    ErrorKind    `json:"kind"`
		Message string       `json:"message"`
		Hint    string       `json:"hint"`
		Source  *file.Source `json:"source"`
		Spans   []ErrorSpan  `json:"spans"`
		Cause   error        `json:"cause"`
	}
)

const (
	UnknownError     ErrorKind = ""
	SyntaxError      ErrorKind = "SyntaxError"
	NameError        ErrorKind = "NameError"
	TypeError        ErrorKind = "TypeError"
	SemanticError    ErrorKind = "SemanticError"
	UnsupportedError ErrorKind = "UnsupportedError"
	InternalError    ErrorKind = "InternalError"
)

func (e *CompilationError) Error() string {
	return e.Message
}

func (e *CompilationError) Format() string {
	var b strings.Builder
	FormatError(&b, e, 0)
	return b.String()
}
