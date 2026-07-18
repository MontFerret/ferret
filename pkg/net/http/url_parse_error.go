package http

// URLParseError reports a failure while parsing a request URL.
type URLParseError struct {
	// Err is the underlying URL parsing failure.
	Err error
}

// Error returns the human-readable URL parsing failure.
func (e *URLParseError) Error() string {
	if e == nil || e.Err == nil {
		return "http: parse url"
	}

	return "http: parse url: " + e.Err.Error()
}

// Unwrap exposes the underlying URL parsing failure.
func (e *URLParseError) Unwrap() error {
	if e == nil {
		return nil
	}

	return e.Err
}
