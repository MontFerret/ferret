package http

import (
	"errors"
	"net/url"
)

// URLParseError reports a failure while parsing a request URL. Its message
// omits the raw URL so malformed credentials are not disclosed; Err retains
// the underlying cause for errors.Is and errors.As.
type URLParseError struct {
	// Err is the underlying URL parsing failure.
	Err error
}

// Error returns the human-readable URL parsing failure.
func (e *URLParseError) Error() string {
	if e == nil || e.Err == nil {
		return "http: parse url"
	}

	var urlErr *url.Error
	if errors.As(e.Err, &urlErr) && urlErr.Err != nil {
		return "http: parse url: " + urlErr.Err.Error()
	}

	return "http: parse url"
}

// Unwrap exposes the underlying URL parsing failure.
func (e *URLParseError) Unwrap() error {
	if e == nil {
		return nil
	}

	return e.Err
}
