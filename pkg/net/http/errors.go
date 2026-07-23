package http

import "errors"

var (
	// ErrNilClient indicates a standard-library HTTP client is nil.
	ErrNilClient = errors.New("http: client is nil")

	// ErrNilRequest indicates an HTTP client received a nil request.
	ErrNilRequest = errors.New("http: request is nil")

	// ErrNilResponse indicates a RoundTripper returned a nil response without
	// reporting an error.
	ErrNilResponse = errors.New("http: response is nil")

	// ErrPolicyDenied identifies requests rejected by the HTTP access policy.
	ErrPolicyDenied = errors.New("http: blocked by access policy")

	// ErrInvalidPolicyConfiguration identifies malformed or contradictory
	// options supplied while constructing a Policy.
	ErrInvalidPolicyConfiguration = errors.New("http: invalid policy configuration")
)
