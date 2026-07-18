package http

import "fmt"

// URLValidationError reports a structurally invalid request URL.
type URLValidationError struct {
	// Field identifies the invalid URL component.
	Field string
	// Reason explains why the component is invalid.
	Reason string
}

// Error returns the human-readable URL validation failure.
func (e *URLValidationError) Error() string {
	if e == nil {
		return "http: invalid url"
	}

	return fmt.Sprintf("http: url %s %s", e.Field, e.Reason)
}
