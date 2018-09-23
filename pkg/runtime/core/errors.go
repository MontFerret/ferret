package core

import (
	"fmt"
	"github.com/pkg/errors"
	"strings"
)

var (
	ErrMissedArgument        = errors.New("missed argument")
	ErrInvalidArgument       = errors.New("invalid argument")
	ErrInvalidArgumentNumber = errors.New("invalid argument number")
	ErrInvalidType           = errors.New("invalid type")
	ErrInvalidOperation      = errors.New("invalid operation")
	ErrNotFound              = errors.New("not found")
	ErrNotUnique             = errors.New("not unique")
	ErrTerminated            = errors.New("operation is terminated")
	ErrUnexpected            = errors.New("unexpected error")
	ErrTimeout               = errors.New("operation timed out")
	ErrNotImplemented        = errors.New("not implemented")
)

const typeErrorTemplate = "expected %s, but got %s"

func SourceError(src SourceMap, err error) error {
	return errors.Errorf("%s: %s", err.Error(), src.String())
}

func TypeError(actual Type, expected ...Type) error {
	if len(expected) == 0 {
		return Error(ErrInvalidType, actual.String())
	}

	if len(expected) == 1 {
		return Error(ErrInvalidType, fmt.Sprintf(typeErrorTemplate, expected, actual))
	}

	strs := make([]string, len(expected))

	for idx, t := range expected {
		strs[idx] = t.String()
	}

	expectedStr := strings.Join(strs, " or ")

	return Error(ErrInvalidType, fmt.Sprintf(typeErrorTemplate, expectedStr, actual))
}

func Error(err error, msg string) error {
	return errors.Errorf("%s: %s", err.Error(), msg)
}
