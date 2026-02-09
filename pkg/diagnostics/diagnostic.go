package diagnostics

import (
	"strings"

	"github.com/MontFerret/ferret/pkg/file"
)

// Diagnostic represents a structured diagnostic error with optional source context and spans.
type Diagnostic struct {
	Kind    Kind         `json:"kind"`
	Message string       `json:"message"`
	Hint    string       `json:"hint"`
	Note    string       `json:"note"`
	Source  *file.Source `json:"source"`
	Spans   []ErrorSpan  `json:"spans"`
	Cause   error        `json:"cause"`
}

// ErrorSpan represents a specific span in the source code related to a diagnostic.
func (e *Diagnostic) String() string {
	return e.Message
}

// Error implements the error interface for Diagnostic, returning the message as the error string.
func (e *Diagnostic) Error() string {
	return e.Message
}

// Unwrap allows for error unwrapping, returning the underlying cause of the diagnostic if it exists.
func (e *Diagnostic) Unwrap() error {
	return e.Cause
}

// Format generates a formatted string representation of the diagnostic, including its message, source context, and spans.
func (e *Diagnostic) Format() string {
	var b strings.Builder

	FormatDiagnostic(&b, e, 0)

	return b.String()
}
