package vm

import (
	"errors"
	"strings"

	"github.com/MontFerret/ferret/pkg/diagnostics"
)

// RuntimeError represents a VM execution error with source context.
type (
	RuntimeError struct {
		*diagnostics.Diagnostic
	}

	warmupError struct {
		Err error
		PC  int
		Dst Operand
	}

	warmupErrorSet struct {
		Errors []*warmupError
	}
)

func (e *warmupError) Error() string {
	return e.Err.Error()
}

func (e *warmupError) Unwrap() error {
	return e.Err
}

func (e *warmupErrorSet) Unwrap() error {
	if len(e.Errors) == 0 {
		return nil
	}

	var err error

	for _, we := range e.Errors {
		if err == nil {
			err = we
		} else {
			err = errors.Join(err, we)
		}
	}

	return err
}

func (e *warmupErrorSet) Error() string {
	if len(e.Errors) == 0 {
		return ""
	}

	var b strings.Builder

	for _, we := range e.Errors {
		b.WriteString(we.Error())
		b.WriteString("\n")
	}

	return b.String()
}

func (e *warmupErrorSet) Add(err error, pc int, dst Operand) {
	e.Errors = append(e.Errors, &warmupError{
		Err: err,
		PC:  pc,
		Dst: dst,
	})
}

func (e *warmupErrorSet) Size() int {
	return len(e.Errors)
}

const (
	ArityError       diagnostics.Kind = "ArityError"
	NullDereferenced diagnostics.Kind = "NullDereference"
	DivideByZero     diagnostics.Kind = "DivideByZero"
	ModuloByZero     diagnostics.Kind = "ModuloByZero"
	AssertionFailed  diagnostics.Kind = "AssertionFailed"
	UnresolvedSymbol diagnostics.Kind = "UnresolvedSymbol"
	UncaughtError    diagnostics.Kind = "UncaughtError"
)

var (
	ErrMissedParam           = errors.New("missed parameter")
	ErrInsufficientRegisters = errors.New("insufficient registers")
	ErrUnresolvedFunction    = errors.New("unresolved function")
	ErrInvalidFunctionName   = errors.New("invalid function name")
	ErrDivisionByZero        = errors.New("division by zero")
	ErrModuloByZero          = errors.New("modulo by zero")
)
