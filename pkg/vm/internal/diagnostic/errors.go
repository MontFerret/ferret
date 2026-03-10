package diagnostic

import (
	"errors"
	"strings"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/diagnostics"
)

// RuntimeError represents a VM execution error with source context.
type RuntimeError struct {
	*diagnostics.Diagnostic
}

type WarmupError struct {
	Err error
	PC  int
	Dst bytecode.Operand
}

type WarmupErrorSet struct {
	Errors []*WarmupError
}

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
	if e == nil {
		return ""
	}

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

func (e *WarmupError) Error() string {
	return e.Err.Error()
}

func (e *WarmupError) Unwrap() error {
	return e.Err
}

func (e *WarmupErrorSet) Unwrap() error {
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

func (e *WarmupErrorSet) Error() string {
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

func (e *WarmupErrorSet) Add(err error, pc int, dst bytecode.Operand) {
	e.Errors = append(e.Errors, &WarmupError{
		Err: err,
		PC:  pc,
		Dst: dst,
	})
}

func (e *WarmupErrorSet) Size() int {
	return len(e.Errors)
}

type RuntimeErrorSet = diagnostics.Diagnostics[*RuntimeError]

func NewRuntimeErrorSet(size int) *RuntimeErrorSet {
	return diagnostics.NewDiagnostics[*RuntimeError](size)
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
