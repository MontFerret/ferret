package runtime

import (
	"fmt"
	"strings"

	"errors"
)

var (
	ErrMissedArgument        = errors.New("missed argument")
	ErrInvalidArgument       = errors.New("invalid argument")
	ErrInvalidArgumentNumber = errors.New("invalid argument number")
	ErrInvalidArgumentType   = errors.New("invalid argument type")
	ErrInvalidType           = errors.New("invalid type")
	ErrInvalidOperation      = errors.New("invalid operation")
	ErrNotFound              = errors.New("not found")
	ErrNotUnique             = errors.New("not unique")
	ErrUnexpected            = errors.New("unexpected error")
	ErrTimeout               = errors.New("operation timed out")
	ErrNotImplemented        = errors.New("not implemented")
	ErrNotSupported          = errors.New("not supported")
	ErrRange                 = errors.New("out of range")
)

const (
	typeErrorTemplate = "expected %s, but got %s"
)

// TypeErrorOf creates a new error indicating that the provided value has an invalid type.
// The expected parameter can be used to specify one or more expected types for the value.
func TypeErrorOf(value Value, expected ...Type) error {
	return TypeError(TypeOf(value), expected...)
}

// TypeError creates a new error indicating that the provided type is invalid.
// The expected parameter can be used to specify one or more expected types for the value.
func TypeError(actual Type, expected ...Type) error {
	if len(expected) == 0 {
		return Error(ErrInvalidType, typeString(actual))
	}

	strs := make([]string, len(expected))

	for idx, t := range expected {
		strs[idx] = typeString(t)
	}

	expectedStr := strings.Join(strs, " or ")

	return Error(ErrInvalidType, fmt.Sprintf(typeErrorTemplate, expectedStr, typeString(actual)))
}

func typeString(t Type) string {
	if t == nil {
		return "Unknown"
	}

	return t.String()
}

// Error creates a new error by wrapping the provided error with the given message.
// The resulting error will include both the original error and the additional message, providing more context about the error.
func Error(err error, msg string) error {
	return fmt.Errorf("%w: %s", err, msg)
}

// Errorf creates a new error by wrapping the provided error with a formatted message.
// The resulting error will include both the original error and the formatted message, providing more context about the error.
func Errorf(err error, format string, args ...interface{}) error {
	return fmt.Errorf("%w: %s", err, fmt.Sprintf(format, args...))
}
