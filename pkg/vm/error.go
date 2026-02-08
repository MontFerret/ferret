package vm

import (
	"fmt"
	"strings"

	shareddiag "github.com/MontFerret/ferret/pkg/diagnostics"
	"github.com/MontFerret/ferret/pkg/file"
)

// RuntimeError represents a VM execution error with source context.
type RuntimeError struct {
	Message string
	Hint    string
	Note    string
	Label   string
	Source  *file.Source
	Span    file.Span
	Cause   error
}

func (e *RuntimeError) Error() string {
	return e.Message
}

func (e *RuntimeError) Unwrap() error {
	return e.Cause
}

func (e *RuntimeError) Format() string {
	var b strings.Builder

	fmt.Fprintf(&b, "error: %s\n", e.Message)

	renderer := shareddiag.SpanRenderer{
		CaretChar:          '^',
		ShowTrailingGutter: true,
	}

	renderer.Render(&b, e.Source, e.Span, e.Label)

	if e.Hint != "" {
		fmt.Fprintf(&b, "Hint: %s\n", e.Hint)
	}

	if e.Note != "" {
		fmt.Fprintf(&b, " = note: %s\n", e.Note)
	}

	if e.Cause != nil && e.Note == "" {
		fmt.Fprintf(&b, " = note: caused by: %s\n", e.Cause.Error())
	}

	return b.String()
}
