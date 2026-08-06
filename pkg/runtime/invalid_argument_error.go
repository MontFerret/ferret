package runtime

import (
	"errors"
	"fmt"
)

type invalidArgumentError struct {
	cause    error
	position int
}

// InvalidArgumentDetails returns the outermost invalid argument position and cause.
// The returned position is zero-based so callers can use it both for metadata lookups
// and for user-facing messages after converting to a one-based index.
func InvalidArgumentDetails(err error) (pos int, ok bool, cause error) {
	var invalidErr *invalidArgumentError

	if !errors.As(err, &invalidErr) || invalidErr == nil {
		return 0, false, nil
	}

	return invalidErr.position, true, invalidErr.cause
}

func newInvalidArgumentError(err error, pos int) error {
	return &invalidArgumentError{
		cause:    err,
		position: pos,
	}
}

func (e *invalidArgumentError) Error() string {
	if e == nil {
		return ErrInvalidArgument.Error()
	}

	if e.cause == nil {
		return fmt.Sprintf("%s at position %d", ErrInvalidArgument, e.position+1)
	}

	return fmt.Sprintf("%s at position %d - %s", ErrInvalidArgument, e.position+1, e.cause)
}

func (e *invalidArgumentError) Unwrap() error {
	if e == nil {
		return nil
	}

	return e.cause
}

func (e *invalidArgumentError) Is(target error) bool {
	return target == ErrInvalidArgument
}
