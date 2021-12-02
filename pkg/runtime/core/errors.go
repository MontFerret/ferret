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
		Cause() error
		Segment() int
		Format(path []Value) string
	}

	// NativePathError represents a default implementation of GetterError interface.
	NativePathError struct {
		cause   error
		segment int
	}
)

// NewPathError is a constructor function of NativePathError struct.
func NewPathError(err error, segment int) PathError {
	return &NativePathError{
		cause:   err,
		segment: segment,
	}
}

// NewPathErrorFrom is a constructor function of NativePathError struct
// that accepts nested PathError and original segment index.
// It sums indexes to get the correct one that points to original path.
func NewPathErrorFrom(err PathError, segment int) PathError {
	return NewPathError(err.Cause(), err.Segment()+segment)
}

func (e *NativePathError) Cause() error {
	return e.cause
}

func (e *NativePathError) Error() string {
	return e.cause.Error()
}

func (e *NativePathError) Segment() int {
	return e.segment
}

func (e *NativePathError) Format(path []Value) string {
	err := e.cause

	if err == ErrInvalidPath && len(path) > e.segment {
		return err.Error() + " '" + path[e.segment].String() + "'"
	}

	return err.Error()
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
	ErrInvalidPath           = errors.New("cannot read property")
	ErrDone                  = errors.New("operation done")
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

func IsDone(err error) bool {
	return err == ErrDone
}
