package vm

import (
	"errors"

	"github.com/MontFerret/ferret/v2/pkg/diagnostics"
	rtdiagnostics "github.com/MontFerret/ferret/v2/pkg/vm/internal/diagnostics"
)

type (
	// RuntimeError represents a VM execution error with source context.
	RuntimeError = rtdiagnostics.RuntimeError

	// WarmupError represents an error that occurs during the warmup phase with details such as PC and destination operand.
	WarmupError = rtdiagnostics.WarmupError

	// InitializationError represents an error occurring during initialization with details like cause, message, PC, and destination operand.
	InitializationError = rtdiagnostics.InitializationError

	// InvariantError represents an invariant violation error with an optional cause and a descriptive message.
	InvariantError = rtdiagnostics.InvariantError
)

type (
	errAction    uint8
	errorKind    uint8
	recoveryMode uint8
)

const (
	errOK errAction = iota
	errContinue
	errReturn
)

const (
	errKindRuntime errorKind = iota
	errKindInvariant
)

const (
	recoverDefault recoveryMode = iota
	recoverOptional
	recoverMissingMember
	recoverProtected
)

const (
	ArityError       diagnostics.Kind = rtdiagnostics.ArityError
	InvalidArgument  diagnostics.Kind = rtdiagnostics.InvalidArgument
	NullDereferenced diagnostics.Kind = rtdiagnostics.NullDereferenced
	DivideByZero     diagnostics.Kind = rtdiagnostics.DivideByZero
	ModuloByZero     diagnostics.Kind = rtdiagnostics.ModuloByZero
	AssertionFailed  diagnostics.Kind = rtdiagnostics.AssertionFailed
	UnresolvedSymbol diagnostics.Kind = rtdiagnostics.UnresolvedSymbol
	UncaughtError    diagnostics.Kind = rtdiagnostics.UncaughtError
)

var (
	ErrMissedParam           = rtdiagnostics.ErrMissedParam
	ErrInsufficientRegisters = rtdiagnostics.ErrInsufficientRegisters
	ErrUnresolvedFunction    = rtdiagnostics.ErrUnresolvedFunction
	ErrInvalidFunctionName   = rtdiagnostics.ErrInvalidFunctionName
	ErrDivisionByZero        = rtdiagnostics.ErrDivisionByZero
	ErrModuloByZero          = rtdiagnostics.ErrModuloByZero
	ErrPoolExhausted         = errors.New("pool exhausted")
	ErrPoolClosed            = errors.New("pool closed")
)
