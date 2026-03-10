package vm

import (
	"github.com/MontFerret/ferret/v2/pkg/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/diagnostic"
)

// RuntimeError represents a VM execution error with source context.
type RuntimeError = diagnostic.RuntimeError

type (
	errAction uint8
	failClass uint8
)

const (
	errOK errAction = iota
	errContinue
	errReturn
)

const (
	failRuntime failClass = iota
	failProtected
	failInvariant
)

const (
	ArityError       diagnostics.Kind = diagnostic.ArityError
	NullDereferenced diagnostics.Kind = diagnostic.NullDereferenced
	DivideByZero     diagnostics.Kind = diagnostic.DivideByZero
	ModuloByZero     diagnostics.Kind = diagnostic.ModuloByZero
	AssertionFailed  diagnostics.Kind = diagnostic.AssertionFailed
	UnresolvedSymbol diagnostics.Kind = diagnostic.UnresolvedSymbol
	UncaughtError    diagnostics.Kind = diagnostic.UncaughtError
)

var (
	ErrMissedParam           = diagnostic.ErrMissedParam
	ErrInsufficientRegisters = diagnostic.ErrInsufficientRegisters
	ErrUnresolvedFunction    = diagnostic.ErrUnresolvedFunction
	ErrInvalidFunctionName   = diagnostic.ErrInvalidFunctionName
	ErrDivisionByZero        = diagnostic.ErrDivisionByZero
	ErrModuloByZero          = diagnostic.ErrModuloByZero
)
