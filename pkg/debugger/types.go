package debugger

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/encoding"
	"github.com/MontFerret/ferret/v2/pkg/source"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

type (
	// Reason identifies why a debug execution stopped.
	Reason string

	// BreakpointID identifies a breakpoint within one debugger session.
	BreakpointID int

	// Location identifies a source location in the debugged file.
	Location struct {
		File   string
		Line   int
		Column int
		Span   source.Span
	}

	// Value is a safely formatted debugger value.
	Value struct {
		Type    string
		Display string
	}

	// Variable describes a visible local or bind parameter.
	Variable struct {
		Name    string
		Value   Value
		Mutable bool
		Param   bool
	}

	// Frame describes the paused top frame or one of its callers.
	Frame struct {
		Name       string
		Location   Location
		FunctionID int
	}

	// Breakpoint describes a requested source-line breakpoint and its resolved
	// executable location, when one exists.
	Breakpoint struct {
		File            string
		RequestedLine   int
		RequestedColumn int
		ID              BreakpointID
		PointID         bytecode.DebugPointID
		FunctionID      int
		Line            int
		Column          int
		Bound           bool
	}

	// Event reports a debugger stop, completion, or termination.
	Event struct {
		Error    error
		Output   *encoding.Output
		Reason   Reason
		Location Location
		Depth    int
	}

	// FormatOptions bounds debugger value traversal and rendered output.
	FormatOptions struct {
		MaxDepth int
		MaxItems int
		MaxBytes int
	}

	// SessionServices supplies embedding-owned lifecycle and output behavior.
	SessionServices interface {
		BeforeRun(context.Context) (context.Context, error)
		AfterRun(context.Context, error) error
		ExtendContext(context.Context) context.Context
		Materialize(*vm.Result) (*encoding.Output, error)
		Close() error
	}

	// Config contains the dependencies for an advanced debugger session.
	Config struct {
		Execution   vm.DebugExecution
		Values      vm.DebugValueAccess
		Services    SessionServices
		Source      *source.Source
		DebugPoints []bytecode.DebugPoint
		Params      []string
		Format      FormatOptions
	}
)

const (
	ReasonEntry        Reason = "entry"
	ReasonBreakpoint   Reason = "breakpoint"
	ReasonStep         Reason = "step"
	ReasonPause        Reason = "pause"
	ReasonRuntimeError Reason = "runtime-error"
	ReasonCompleted    Reason = "completed"
	ReasonTerminated   Reason = "terminated"
)

// DefaultFormatOptions returns conservative debugger formatting limits.
func DefaultFormatOptions() FormatOptions {
	return FormatOptions{MaxDepth: 3, MaxItems: 8, MaxBytes: 1024}
}
