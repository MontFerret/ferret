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
	ErrTerminated            = errors.New("operation is terminated")
	ErrUnexpected            = errors.New("unexpected error")
	ErrTimeout               = errors.New("operation timed out")
	ErrNotImplemented        = errors.New("not implemented")
	ErrNotSupported          = errors.New("not supported")
)

const typeErrorTemplate = "expected %s, but got %s"

func TypeErrorOf(value Value, expected ...Type) error {
	return TypeError(Reflect(value), expected...)
}

func TypeError(actual Type, expected ...Type) error {
	if len(expected) == 0 {
		return Error(ErrInvalidType, string(actual))
	}

	if len(expected) == 1 {
		return Error(ErrInvalidType, fmt.Sprintf(typeErrorTemplate, expected, actual))
	}

	strs := make([]string, len(expected))

	for idx, t := range expected {
		strs[idx] = string(t)
	}

	expectedStr := strings.Join(strs, " or ")

	return Error(ErrInvalidType, fmt.Sprintf(typeErrorTemplate, expectedStr, actual))
}

func Error(err error, msg string) error {
	return fmt.Errorf("%w: %s", err, msg)
}

func Errorf(err error, format string, args ...interface{}) error {
	return fmt.Errorf("%w: %s", err, fmt.Sprintf(format, args...))
}

func Errors(err ...error) error {
	message := ""

	for _, e := range err {
		if e != nil {
			message += ": " + e.Error()
		}
	}

	return errors.New(message)
}

func Errorsf(msg string, err ...error) error {
	if len(err) == 0 {
		return errors.New(msg)
	}

	res := err[0]

	for i := 1; i < len(err); i++ {
		e := err[i]
		if e != nil {
			res = fmt.Errorf("%w: %s", e, msg)
		}
	}

	return fmt.Errorf("%s: %w", msg, res)
}
