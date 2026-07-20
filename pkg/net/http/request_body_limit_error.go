package http

import "fmt"

// RequestBodyLimitError reports a request body larger than the configured limit.
type RequestBodyLimitError struct {
	// Size is the request body size in bytes.
	Size int64
	// Limit is the configured request body limit in bytes.
	Limit int64
}

// Error returns the human-readable request body limit failure.
func (e *RequestBodyLimitError) Error() string {
	if e == nil {
		return "http: request body exceeds limit"
	}

	return fmt.Sprintf("http: request body exceeds limit: %d > %d", e.Size, e.Limit)
}
