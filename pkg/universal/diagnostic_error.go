package universal

import apidiagnostics "github.com/MontFerret/api/diagnostics"

type diagnosticError struct {
	cause       error
	diagnostics apidiagnostics.Diagnostics
}

func (e *diagnosticError) Error() string {
	return e.cause.Error()
}

func (e *diagnosticError) Unwrap() []error {
	// Expose the complete projection before a possibly already-projected cause,
	// so errors.As cannot stop at an older, incomplete diagnostic collection.
	return []error{e.diagnostics, e.cause}
}

func newDiagnosticError(cause error, diagnostics apidiagnostics.Diagnostics) error {
	return &diagnosticError{cause: cause, diagnostics: diagnostics}
}
