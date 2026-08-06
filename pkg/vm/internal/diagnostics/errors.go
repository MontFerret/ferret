package diagnostics

import (
	"errors"

	"github.com/MontFerret/ferret/v2/pkg/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/internal/operationerror"
)

const (
	ArityError       diagnostics.Kind = "ArityError"
	InvalidArgument  diagnostics.Kind = "InvalidArgument"
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
	ErrDivisionByZero        = operationerror.ErrDivisionByZero
	ErrModuloByZero          = operationerror.ErrModuloByZero
)
