package runtime

// conversionError distinguishes failures produced by conversion itself from
// operational errors returned by fallible runtime values.
type conversionError struct {
	target Type
	cause  error
}

func newConversionError(target Type, cause error) error {
	return &conversionError{
		target: target,
		cause:  cause,
	}
}

func (e *conversionError) Error() string {
	return e.cause.Error()
}

func (e *conversionError) Unwrap() error {
	return e.cause
}

func (e *conversionError) targets(target Type) bool {
	return e.target == target
}
