package runtime

import "fmt"

type invalidArgumentError struct {
	cause    error
	position int
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
