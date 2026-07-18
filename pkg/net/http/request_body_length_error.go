package http

import "fmt"

// RequestBodyLengthError reports a non-empty request body whose length is
// unknown while a finite request-body limit is configured.
type RequestBodyLengthError struct {
	// ContentLength is the request's reported content length.
	ContentLength int64
	// Limit is the configured request body limit in bytes.
	Limit int64
}

// Error returns the human-readable unknown request-body length failure.
func (e *RequestBodyLengthError) Error() string {
	if e == nil {
		return "http: request body length is unknown under finite limit"
	}

	return fmt.Sprintf(
		"http: request body length is unknown under finite limit: content-length=%d limit=%d",
		e.ContentLength,
		e.Limit,
	)
}
