package vm

import (
	"errors"

	"github.com/MontFerret/ferret/pkg/diagnostics"
)

// RuntimeError represents a VM execution error with source context.
type RuntimeError struct {
	*diagnostics.Diagnostic
}

const (
	ArityError       diagnostics.Kind = "ArityError"
	NullDereferenced diagnostics.Kind = "NullDereference"
	DivideByZero     diagnostics.Kind = "DivideByZero"
	ModuloByZero     diagnostics.Kind = "ModuloByZero"
	AssertionFailed  diagnostics.Kind = "AssertionFailed"
	UnresolvedSymbol diagnostics.Kind = "UnresolvedSymbol"
)

var (
	ErrMissedParam           = errors.New("missed parameter")
	ErrFunctionNotFound      = errors.New("function not found")
	ErrInsufficientRegisters = errors.New("insufficient registers")
	ErrUnresolvedFunction    = errors.New("unresolved function")
	ErrInvalidFunctionName   = errors.New("invalid function name")
	ErrDivisionByZero        = errors.New("division by zero")
	ErrModuloByZero          = errors.New("modulo by zero")
)
