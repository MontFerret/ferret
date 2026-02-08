package diagnostics

import (
	"fmt"
	"strings"
)

type MultiError struct {
	errors []error
}

func NewMultiError(errors ...error) *MultiError {
	return &MultiError{
		errors: errors,
	}
}

func (e *MultiError) Errors() []error {
	return e.errors
}

func (e *MultiError) Error() string {
	if len(e.errors) == 0 {
		return "No errors"
	}

	return fmt.Sprintf("Found %d errors", len(e.errors))
}

func (e *MultiError) Format() string {
	if len(e.errors) == 0 {
		return ""
	}

	var b strings.Builder

	for _, err := range e.errors {
		if err != nil {
			if formattable, ok := err.(Formattable); ok {
				b.WriteString(formattable.Format())
			} else {
				b.WriteString(err.Error())
			}

			b.WriteString("\n")
		}
	}

	return b.String()
}
