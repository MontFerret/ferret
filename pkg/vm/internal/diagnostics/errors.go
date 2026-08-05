package diagnostics

import (
	"errors"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
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
	ErrDivisionByZero        = errors.New("division by zero")
	ErrModuloByZero          = errors.New("modulo by zero")
)

func CheckDivisionByZero(
	program *bytecode.Program,
	pc int,
	left runtime.Value,
	right runtime.Value,
) error {
	if validDivisionPair(left, right) && isZeroDivisor(right) {
		return NewRuntimeError(
			program,
			pc,
			DivideByZero,
			"division by zero",
			"attempt to divide by zero",
			"Ensure the denominator is non-zero before division",
			"",
			ErrDivisionByZero,
		)
	}

	return nil
}

func CheckModuloByZero(
	program *bytecode.Program,
	pc int,
	left runtime.Value,
	right runtime.Value,
) error {
	if isNativeNumber(left) && isNativeNumber(right) && isZeroDivisor(right) {
		return NewRuntimeError(
			program,
			pc,
			ModuloByZero,
			"modulo by zero",
			"attempt to take modulo by zero",
			"Ensure the divisor is non-zero before modulo",
			"",
			ErrModuloByZero,
		)
	}

	return nil
}

func validDivisionPair(left, right runtime.Value) bool {
	if isNativeNumber(left) && isNativeNumber(right) {
		return true
	}

	if _, ok := left.(runtime.Duration); !ok {
		return false
	}

	switch right.(type) {
	case runtime.Int, runtime.Float, runtime.Duration:
		return true
	default:
		return false
	}
}

func isNativeNumber(value runtime.Value) bool {
	switch value.(type) {
	case runtime.Int, runtime.Float:
		return true
	default:
		return false
	}
}

func isZeroDivisor(value runtime.Value) bool {
	switch value := value.(type) {
	case runtime.Int:
		return value == 0
	case runtime.Float:
		return value == 0
	case runtime.Duration:
		return value == 0
	default:
		return false
	}
}
