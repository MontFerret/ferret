package core

import (
	"strings"

	"github.com/MontFerret/ferret/pkg/file"
)

type (
	ErrorKind string

	CompilationError struct {
		Message  string         `json:"message"`
		Kind     ErrorKind      `json:"kind"`
		Source   *file.Source   `json:"source"`
		Location *file.Location `json:"location"`
		Hint     string         `json:"hint"`
		Cause    error          `json:"cause"`
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
