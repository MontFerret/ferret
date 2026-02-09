package diagnostics

import (
	"fmt"
	"strings"
)

type Diagnostics[E FormattableError] struct {
	errors []E
}

func NewDiagnostics[E FormattableError](errors []E) *Diagnostics[E] {
	return &Diagnostics[E]{
		errors: errors,
	}
}

func (e *Diagnostics[E]) Size() int {
	return len(e.errors)
}

func (e *Diagnostics[E]) Add(err E) {
	e.errors = append(e.errors, err)
}

func (e *Diagnostics[E]) Errors() []E {
	return e.errors
}

func (e *Diagnostics[E]) Error() string {
	if len(e.errors) == 0 {
		return "No errors"
	}

	return fmt.Sprintf("Found %d errors", len(e.errors))
}

func (e *Diagnostics[E]) Format() string {
	if len(e.errors) == 0 {
		return ""
	}

	var b strings.Builder

	for _, err := range e.errors {
		b.WriteString(err.Format())
		b.WriteString("\n")
	}

	return b.String()
}
