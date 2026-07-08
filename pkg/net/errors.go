package net

import "errors"

var (
	// ErrNotFound indicates no Network is available in the provided context.
	ErrNotFound = errors.New("network not found in context")

	// ErrHTTPClientNotFound indicates no HTTP client is available in the provided context.
	ErrHTTPClientNotFound = errors.New("http client not found in context")
)
