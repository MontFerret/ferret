package debugger

import (
	"context"

	apidebugger "github.com/MontFerret/api/debugger"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/encoding"
	"github.com/MontFerret/ferret/v2/pkg/source"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

type (
	// Reason identifies why execution stopped.
	Reason = apidebugger.Reason
	// BreakpointID identifies a breakpoint within one session.
	BreakpointID = apidebugger.BreakpointID
	// ValueReference identifies an expandable value in one paused state.
	ValueReference = apidebugger.ValueReference
	// BreakpointBindingMode selects how a source location binds.
	BreakpointBindingMode = apidebugger.BreakpointBindingMode
	// Value is a safely formatted portable debugger value.
	Value = apidebugger.Value
	// Variable describes a local or bind parameter.
	Variable = apidebugger.Variable
	// Frame describes the paused frame or one of its callers.
	Frame = apidebugger.Frame
	// Breakpoint describes requested and resolved source locations.
	Breakpoint = apidebugger.Breakpoint
	// BreakpointOptions configures source-location binding.
	BreakpointOptions = apidebugger.BreakpointOptions
	// Event reports a debugger stop, completion, or termination.
	Event = apidebugger.Event

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
		DebugPoints []bytecode.DebugPoint
		Params      []string
		Source      source.Source
		Format      FormatOptions
	}
)

const (
	ReasonEntry                            = apidebugger.ReasonEntry
	ReasonBreakpoint                       = apidebugger.ReasonBreakpoint
	ReasonStep                             = apidebugger.ReasonStep
	ReasonPause                            = apidebugger.ReasonPause
	ReasonRuntimeError                     = apidebugger.ReasonRuntimeError
	ReasonCompleted                        = apidebugger.ReasonCompleted
	ReasonTerminated                       = apidebugger.ReasonTerminated
	BreakpointBindNextExecutableInSource   = apidebugger.BreakpointBindNextExecutableInSource
	BreakpointBindExact                    = apidebugger.BreakpointBindExact
	BreakpointBindNextExecutableInFunction = apidebugger.BreakpointBindNextExecutableInFunction
)

// DefaultFormatOptions returns conservative debugger formatting limits.
func DefaultFormatOptions() FormatOptions {
	return FormatOptions{MaxDepth: 3, MaxItems: 8, MaxBytes: 1024}
}

var _ apidebugger.Session = (*Session)(nil)
