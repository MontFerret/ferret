package http

import "fmt"

// ResponseBodyLimitError reports a response body larger than the configured limit.
type ResponseBodyLimitError struct {
	// Size is the number of response body bytes observed before reading stopped.
	Size int64
	// Limit is the configured response body limit in bytes.
	Limit int64
}

// Error returns the human-readable response body limit failure.
func (e *ResponseBodyLimitError) Error() string {
	if e == nil {
		return "http: response body exceeds limit"
	}

	return fmt.Sprintf("http: response body exceeds limit: %d > %d", e.Size, e.Limit)
}
