package core

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

type (
	SourceErrorDetail struct {
		error
		BaseError    error
		ComputeError error
	}

	// PathError represents an interface of
	// error type which returned when an error occurs during an execution of Getter.GetIn or Setter.SetIn functions
	// and contains segment of a given path that caused the error.
	PathError interface {
		error
		Path() []Value
		Segment() int64
	}

	// DefaultPathError represents a default implementation of GetterError interface.
	DefaultPathError struct {
		cause   error
		path    []Value
		segment int64
	}
)

// NewPathError is a constructor function of DefaultPathError struct.
func NewPathError(err error, path []Value, segment int64) PathError {
	return &DefaultPathError{
		cause:   err,
		path:    path,
		segment: segment,
	}
}

func (e *DefaultPathError) Error() string {
	return e.cause.Error()
}

func (e *DefaultPathError) Path() []Value {
	return e.path
}

func (e *DefaultPathError) Segment() int64 {
	return e.segment
}

func (e *SourceErrorDetail) Error() string {
	return e.ComputeError.Error()
}

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
	ErrNotSupported          = errors.New("not supported")
	ErrNoMoreData            = errors.New("no more data")
)

const typeErrorTemplate = "expected %s, but got %s"

func SourceError(src SourceMap, err error) error {
	return &SourceErrorDetail{
		BaseError:    err,
		ComputeError: errors.Errorf("%s: %s", err.Error(), src.String()),
	}
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

func Errorf(err error, format string, args ...interface{}) error {
	return errors.Errorf("%s: %s", err.Error(), fmt.Sprintf(format, args...))
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

func IsNoMoreData(err error) bool {
	return err == ErrNoMoreData
}
