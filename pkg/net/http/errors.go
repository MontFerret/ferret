package http

import "errors"

// ErrNilRequest indicates an HTTP client received a nil request.
var ErrNilRequest = errors.New("http: request is nil")
