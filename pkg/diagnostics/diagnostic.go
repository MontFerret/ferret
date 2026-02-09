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

func (e *Diagnostic) String() string {
	return e.Message
}

func (e *Diagnostic) Error() string {
	return e.Message
}

func (e *Diagnostic) Unwrap() error {
	return e.Cause
}

func (e *Diagnostic) Format() string {
	var b strings.Builder

	FormatDiagnostic(&b, e, 0)

	return b.String()
}
