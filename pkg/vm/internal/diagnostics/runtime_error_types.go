package diagnostics

import "github.com/MontFerret/ferret/v2/pkg/diagnostics"

type (
	// RuntimeError represents a VM execution error with source context.
	RuntimeError struct {
		*diagnostics.Diagnostic
	}

	// RuntimeErrorSet is a specialized diagnostics.Diagnostics type for RuntimeError.
	RuntimeErrorSet struct {
		diagnostics.Diagnostics[*RuntimeError]
	}
)

// Unwrap exposes the source diagnostic and, through it, the original cause.
func (e *RuntimeError) Unwrap() error {
	if e == nil || e.Diagnostic == nil {
		return nil
	}

	return e.Diagnostic
}
