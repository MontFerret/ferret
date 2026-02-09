package diagnostics

import (
	"fmt"
	"strings"
)

type Diagnostics[E FormattableError] struct {
	errors []E
}

// NewDiagnostics creates a new Diagnostics instance with a preallocated slice for storing errors of the specified size.
func NewDiagnostics[E FormattableError](size int) *Diagnostics[E] {
	return &Diagnostics[E]{
		errors: make([]E, 0, size),
	}
}

// NewDiagnosticsOf creates a new Diagnostics instance with the provided errors.
func NewDiagnosticsOf[E FormattableError](errors []E) *Diagnostics[E] {
	return &Diagnostics[E]{
		errors: errors,
	}
}

// Size returns the number of errors currently stored in the Diagnostics instance.
func (e *Diagnostics[E]) Size() int {
	return len(e.errors)
}

// Add appends a new error to the Diagnostics instance.
func (e *Diagnostics[E]) Add(err E) {
	e.errors = append(e.errors, err)
}

// Errors returns a slice of all errors currently stored in the Diagnostics instance.
func (e *Diagnostics[E]) Errors() []E {
	return e.errors
}

// Error implements the error interface for Diagnostics, returning a summary of the number of errors contained within the Diagnostics instance.
func (e *Diagnostics[E]) Error() string {
	if len(e.errors) == 0 {
		return "No errors"
	}

	return fmt.Sprintf("Found %d errors", len(e.errors))
}

// Format returns a formatted string representation of all errors contained within the Diagnostics instance, with each error on a new line.
// If there are no errors, it returns an empty string.
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
