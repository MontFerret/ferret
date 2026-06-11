package ferret

import (
	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/encoding"
	"github.com/MontFerret/ferret/v2/pkg/fs"
	"github.com/MontFerret/ferret/v2/pkg/logging"
	"github.com/MontFerret/ferret/v2/pkg/source"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

// DebugReason identifies why a debug execution stopped.
type DebugReason string

const (
	DebugReasonEntry        DebugReason = "entry"
	DebugReasonBreakpoint   DebugReason = "breakpoint"
	DebugReasonStep         DebugReason = "step"
	DebugReasonPause        DebugReason = "pause"
	DebugReasonRuntimeError DebugReason = "runtime-error"
	DebugReasonCompleted    DebugReason = "completed"
	DebugReasonTerminated   DebugReason = "terminated"
)

type (
	// DebugLocation identifies a source location in the debugged file.
	DebugLocation struct {
		File   string
		Line   int
		Column int
		Span   source.Span
	}

	// DebugValue is a safely formatted debugger value.
	DebugValue struct {
		Type    string
		Display string
	}

	// DebugVariable describes a visible local or bind parameter.
	DebugVariable struct {
		Name    string
		Value   DebugValue
		Mutable bool
		Param   bool
	}

	// DebugFrame describes the paused top frame or one of its callers.
	DebugFrame struct {
		Name       string
		Location   DebugLocation
		FunctionID int
	}

	// DebugBreakpoint describes a requested source-line breakpoint and its
	// resolved executable location, when one exists.
	DebugBreakpoint struct {
		File          string
		RequestedLine int
		ID            int
		Line          int
		Column        int
		Bound         bool
	}

	// DebugEvent reports a debugger stop, completion, or termination.
	DebugEvent struct {
		Error    error
		Output   *Output
		Reason   DebugReason
		Location DebugLocation
		Depth    int
	}

	debugSessionConfig struct {
		logger            logging.Logger
		hooks             sessionHooks
		fs                fs.FileSystem
		execution         *vm.DebugExecution
		source            *source.Source
		limiter           *sessionLimiter
		encoding          *encoding.Registry
		outputContentType string
		debugPoints       []bytecode.DebugPoint
		params            []string
		format            vm.DebugFormatOptions
	}
)
