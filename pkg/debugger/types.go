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

	// ValueReference identifies an expandable debugger value within one paused
	// session state. References are invalidated when execution starts or resumes.
	ValueReference int

	// BreakpointBindingMode selects how a requested source location resolves to
	// an executable debug point.
	BreakpointBindingMode int

	// Value is a safely formatted debugger value.
	Value struct {
		Type      string
		Display   string
		Reference ValueReference
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
		Location   source.Location
		FunctionID bytecode.FunctionID
	}

	// Breakpoint describes a requested source-location breakpoint and its resolved
	// executable location, when one exists.
	Breakpoint struct {
		RequestedLocation source.Location `json:"requestedLocation"`
		Location          source.Range    `json:"location"`
		ID                BreakpointID
		PointID           bytecode.DebugPointID
		FunctionID        bytecode.FunctionID
		BindingMode       BreakpointBindingMode
		Bound             bool
	}

	// BreakpointOptions configures how a requested source location binds.
	BreakpointOptions struct {
		BindingMode BreakpointBindingMode
	}

	// Event reports a debugger stop, completion, or termination.
	Event struct {
		Error            error
		Output           *encoding.Output
		Reason           Reason
		HitBreakpointIDs []BreakpointID
		Location         source.Range
		Depth            int
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
		DebugPoints []bytecode.DebugPoint
		Params      []string
		Source      source.Source
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

const (
	// BreakpointBindNextExecutableInFile preserves the friendly legacy binding
	// behavior and is the zero-value default.
	BreakpointBindNextExecutableInFile BreakpointBindingMode = iota
	BreakpointBindExact
	BreakpointBindNextExecutableInFunction
)

// DefaultFormatOptions returns conservative debugger formatting limits.
func DefaultFormatOptions() FormatOptions {
	return FormatOptions{MaxDepth: 3, MaxItems: 8, MaxBytes: 1024}
}

// Valid reports whether the reference can be used to request child variables.
func (r ValueReference) Valid() bool {
	return r > 0
}
