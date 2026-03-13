package diagnostics

import (
	"context"
	"errors"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

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

func CheckDivisionByZero(
	ctx context.Context,
	program *bytecode.Program,
	pc int,
	left runtime.Value,
	right runtime.Value,
) error {
	l := runtime.ToNumberOnly(ctx, left)
	if _, ok := l.(runtime.Int); !ok {
		return nil
	}

	r := runtime.ToNumberOnly(ctx, right)
	if rv, ok := r.(runtime.Int); ok && rv == 0 {
		return NewRuntimeError(
			program,
			pc,
			DivideByZero,
			"Division by zero",
			"attempt to divide by zero",
			"Ensure the denominator is non-zero before division",
			"Add a conditional check before dividing",
		)
	}

	return nil
}

func CheckModuloByZero(
	ctx context.Context,
	program *bytecode.Program,
	pc int,
	right runtime.Value,
) error {
	rv, err := runtime.ToInt(ctx, right)
	if err != nil {
		// Keep modulo diagnostics type-safe for invalid string inputs like "x".
		if _, ok := right.(runtime.String); ok {
			return runtime.TypeErrorOf(right, runtime.TypeInt)
		}

		return err
	}

	if rv == 0 {
		return NewRuntimeError(
			program,
			pc,
			ModuloByZero,
			"Modulo by zero",
			"attempt to take modulo by zero",
			"Ensure the divisor is non-zero before modulo",
			"Add a conditional check before modulo",
		)
	}

	return nil
}
