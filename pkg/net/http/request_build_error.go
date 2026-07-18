package http

// RequestBuildError reports a failure while converting a Ferret request into
// a standard-library HTTP request.
type RequestBuildError struct {
	// Err is the underlying request construction failure.
	Err error
}

// Error returns the human-readable request construction failure.
func (e *RequestBuildError) Error() string {
	if e == nil || e.Err == nil {
		return "http: build request"
	}

	return "http: build request: " + e.Err.Error()
}

// Unwrap exposes the underlying request construction failure.
func (e *RequestBuildError) Unwrap() error {
	if e == nil {
		return nil
	}

	return e.Err
}
