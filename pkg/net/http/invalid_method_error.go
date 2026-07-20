package http

import "fmt"

// InvalidMethodError reports a request method that is not a valid HTTP token.
type InvalidMethodError struct {
	// Method is the rejected, unnormalized request method.
	Method string
}

// Error returns the human-readable invalid-method failure.
func (e *InvalidMethodError) Error() string {
	if e == nil {
		return "http: invalid method"
	}

	return fmt.Sprintf("http: invalid method %q", e.Method)
}
