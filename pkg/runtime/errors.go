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

const (
	typeErrorTemplate  = "expected %s, but got %s"
	proxyErrorTemplate = "proxy target %T does not implement %s interface"
)

func TypeErrorOf(value Value, expected ...Type) error {
	return TypeError(TypeOf(value), expected...)
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

func ProxyError(target any, expected ...Type) error {
	if len(expected) == 0 {
		return Error(ErrInvalidType, fmt.Sprintf(proxyErrorTemplate, target, "unknown"))
	}

	if len(expected) == 1 {
		return Error(ErrInvalidType, fmt.Sprintf(proxyErrorTemplate, target, expected[0]))
	}

	strs := make([]string, len(expected))

	for idx, t := range expected {
		strs[idx] = string(t)
	}

	expectedStr := strings.Join(strs, " or ")

	return Error(ErrInvalidType, fmt.Sprintf(proxyErrorTemplate, target, expectedStr))
}

func Error(err error, msg string) error {
	return fmt.Errorf("%w: %s", err, msg)
}

func Errorf(err error, format string, args ...interface{}) error {
	return fmt.Errorf("%w: %s", err, fmt.Sprintf(format, args...))
}

func Errors(errs ...error) error {
	if len(errs) == 0 {
		return nil
	}

	err := errs[0]

	for _, e := range errs[1:] {
		if e != nil {
			err = fmt.Errorf("%w: %w", e, err)
		}
	}

	return err
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
