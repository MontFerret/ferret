package http

import "errors"

var (
	// ErrNilRequest indicates an HTTP client received a nil request.
	ErrNilRequest = errors.New("http: request is nil")

	// ErrPolicyDenied identifies requests rejected by the HTTP access policy.
	ErrPolicyDenied = errors.New("http: blocked by access policy")
)
