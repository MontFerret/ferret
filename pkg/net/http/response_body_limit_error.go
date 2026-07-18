package http

import "fmt"

// ResponseBodyLimitError reports a response body larger than the configured limit.
type ResponseBodyLimitError struct {
	// Limit is the configured response body limit in bytes.
	Limit int64
}

// Error returns the human-readable response body limit failure.
func (e *ResponseBodyLimitError) Error() string {
	if e == nil {
		return "http: response body exceeds limit"
	}

	return fmt.Sprintf("http: response body exceeds limit of %d bytes", e.Limit)
}
