package runtime

// durationConversionError distinguishes failures produced by conversion itself
// from operational errors returned by fallible List implementations.
type durationConversionError struct {
	cause error
}

func newDurationConversionError(cause error) error {
	return &durationConversionError{cause: cause}
}

func (e *durationConversionError) Error() string {
	return e.cause.Error()
}

func (e *durationConversionError) Unwrap() error {
	return e.cause
}
