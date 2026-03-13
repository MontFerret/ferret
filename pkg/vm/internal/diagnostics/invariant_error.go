package diagnostics

import "fmt"

type InvariantError struct {
	Cause   error
	Message string
}

func NewInvariantError(message string, cause error) error {
	return &InvariantError{
		Message: message,
		Cause:   cause,
	}
}

func (e *InvariantError) Error() string {
	if e.Cause != nil {
		return e.Message + ": " + e.Cause.Error()
	}

	return e.Message
}

func (e *InvariantError) Unwrap() error {
	if e == nil {
		return nil
	}

	return e.Cause
}

func (e *InvariantError) Format() string {
	return fmt.Sprintf("invariant error: %s", e.Message)
}
