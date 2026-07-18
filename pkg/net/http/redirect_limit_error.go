package http

import "fmt"

// RedirectLimitError reports that a request exceeded its redirect limit.
type RedirectLimitError struct {
	// Limit is the maximum number of redirects followed before the request stopped.
	Limit int
}

// Error returns the human-readable redirect limit failure.
func (e *RedirectLimitError) Error() string {
	if e == nil {
		return "http: redirect limit exceeded"
	}

	return fmt.Sprintf("http: stopped after %d redirect(s)", e.Limit)
}
