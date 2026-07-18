package http

import "fmt"

// HeaderValidationError reports a malformed request header name or value.
// Header values are deliberately omitted so errors cannot disclose secrets.
type HeaderValidationError struct {
	// Header is the rejected header name.
	Header string
	// Reason explains why the name or one of its values is invalid.
	Reason string
}

// Error returns the human-readable header validation failure.
func (e *HeaderValidationError) Error() string {
	if e == nil {
		return "http: invalid request header"
	}

	return fmt.Sprintf("http: invalid request header %q: %s", e.Header, e.Reason)
}
