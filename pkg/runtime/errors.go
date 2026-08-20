package runtime

import (
	"fmt"
	"strings"

	"errors"

	"github.com/MontFerret/ferret/v2/pkg/internal/operator"
)

var (
	ErrMissedArgument        = errors.New("missed argument")
	ErrInvalidArgument       = errors.New("invalid argument")
	ErrInvalidArgumentNumber = errors.New("invalid argument number")
	ErrInvalidArgumentType   = errors.New("invalid argument type")
	ErrInvalidType           = errors.New("invalid type")
	ErrInvalidOperation      = errors.New("invalid operation")
	// ErrUnsupportedOperands lets a host arithmetic capability decline an operand pair.
	// Runtime operator dispatch continues with the reflected right-hand method and
	// converts a fully declined operation into the normal ErrInvalidOperation error.
	// Implementations may return it directly or wrap it for errors.Is compatibility.
	// All other errors are treated as execution failures and propagate immediately.
	ErrUnsupportedOperands = errors.New("unsupported operands")
	// ErrDivisionByZero identifies an arithmetic division with a zero divisor.
	ErrDivisionByZero = errors.New("division by zero")
	// ErrModuloByZero identifies an arithmetic modulo operation with a zero divisor.
	ErrModuloByZero   = errors.New("modulo by zero")
	ErrNotFound       = errors.New("not found")
	ErrNotUnique      = errors.New("not unique")
	ErrUnexpected     = errors.New("unexpected error")
	ErrTimeout        = errors.New("operation timed out")
	ErrNotImplemented = errors.New("not implemented")
	ErrNotSupported   = errors.New("not supported")
	ErrRange          = errors.New("out of range")
)

const (
	typeErrorTemplate = "expected %s, but got %s"
)

// TypeErrorOf creates a new error indicating that the provided value has an invalid type.
// The expected parameter can be used to specify one or more expected types for the value.
func TypeErrorOf(value Value, expected ...Type) error {
	return TypeError(TypeOf(value), expected...)
}

// TypeError creates a new error indicating that the provided type is invalid.
// The expected parameter can be used to specify one or more expected types for the value.
func TypeError(actual Type, expected ...Type) error {
	if len(expected) == 0 {
		return Error(ErrInvalidType, typeString(actual))
	}

	strs := make([]string, len(expected))

	for idx, t := range expected {
		strs[idx] = typeString(t)
	}

	expectedStr := strings.Join(strs, " or ")

	return Error(ErrInvalidType, fmt.Sprintf(typeErrorTemplate, expectedStr, typeString(actual)))
}

// Error creates a new error by wrapping the provided error with the given message.
// The resulting error will include both the original error and the additional message, providing more context about the error.
func Error(err error, msg string) error {
	return fmt.Errorf("%w: %s", err, msg)
}

// Errorf creates a new error by wrapping the provided error with a formatted message.
// The resulting error will include both the original error and the formatted message, providing more context about the error.
func Errorf(err error, format string, args ...any) error {
	return fmt.Errorf("%w: %s", err, fmt.Sprintf(format, args...))
}

func binaryOperatorTypeError(op operator.Binary, left, right Value) error {
	return Error(
		ErrInvalidOperation,
		operator.CannotApply(op, TypeName(TypeOf(left)), TypeName(TypeOf(right))),
	)
}

func unaryOperatorTypeError(op operator.Unary, value Value) error {
	return Error(
		ErrInvalidOperation,
		operator.CannotApplyUnary(op, TypeName(TypeOf(value))),
	)
}

func incompatibleComparisonError(left, right Value) error {
	return Errorf(
		ErrInvalidOperation,
		"comparison cannot be applied to %s and %s",
		TypeName(TypeOf(left)),
		TypeName(TypeOf(right)),
	)
}

func divisionByZeroError() error {
	return fmt.Errorf("%w: %w", ErrInvalidOperation, ErrDivisionByZero)
}

func moduloByZeroError() error {
	return fmt.Errorf("%w: %w", ErrInvalidOperation, ErrModuloByZero)
}

func typeString(t Type) string {
	if t == nil {
		return "Unknown"
	}

	return t.String()
}
