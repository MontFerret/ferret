package core

import (
	"fmt"
	"strings"
)

type MultiCompilationError struct {
	Errors []*CompilationError
}

func NewMultiCompilationError(errors []*CompilationError) error {
	return &MultiCompilationError{
		Errors: errors,
	}
}

func (e *MultiCompilationError) Error() string {
	if len(e.Errors) == 0 {
		return "No errors"
	}

	return fmt.Sprintf("Found %d errors", len(e.Errors))
}

func (e *MultiCompilationError) Format() string {
	if len(e.Errors) == 0 {
		return "No errors"
	}

	var b strings.Builder

	for i, err := range e.Errors {
		if i > 0 {
			b.WriteString("\n")
		}

		FormatError(&b, err, 0)
	}

	return b.String()
}
