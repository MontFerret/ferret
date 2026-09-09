package ferret

import (
	"github.com/MontFerret/ferret/v2/pkg/debugger"
	"github.com/MontFerret/ferret/v2/pkg/source"
)

type (
	// DebugSession retains execution state for source-level debugging.
	DebugSession = debugger.Session

	// DebugReason identifies why a debug session stopped.
	DebugReason = debugger.Reason

	// DebugLocation identifies the source range of a debug stop.
	DebugLocation = source.Range

	// DebugValue describes a value inspected in a paused debug session.
	DebugValue = debugger.Value

	// DebugVariable describes a visible variable in a paused debug session.
	DebugVariable = debugger.Variable

	// DebugFrame describes a retained call frame.
	DebugFrame = debugger.Frame

	// DebugBreakpoint describes a requested breakpoint and its executable binding.
	DebugBreakpoint = debugger.Breakpoint

	// DebugBreakpointID identifies a breakpoint within its debug session.
	DebugBreakpointID = debugger.BreakpointID

	// DebugBreakpointOptions configures how a requested source location binds.
	DebugBreakpointOptions = debugger.BreakpointOptions

	// DebugBreakpointBindingMode selects how a requested source location is bound.
	DebugBreakpointBindingMode = debugger.BreakpointBindingMode

	// DebugSourceLocation identifies a position requested for a breakpoint.
	DebugSourceLocation = source.Location

	// DebugValueReference identifies an expandable value in the current paused state.
	DebugValueReference = debugger.ValueReference

	// DebugEvent describes a debug stop, failure, or completed output.
	DebugEvent = debugger.Event

	// DebugStateError reports an operation that is invalid in the session's current state.
	DebugStateError = debugger.StateError

	// DebugFormatOptions bounds debugger value formatting.
	DebugFormatOptions = debugger.FormatOptions
)

const (
	DebugReasonEntry        = debugger.ReasonEntry
	DebugReasonBreakpoint   = debugger.ReasonBreakpoint
	DebugReasonStep         = debugger.ReasonStep
	DebugReasonPause        = debugger.ReasonPause
	DebugReasonRuntimeError = debugger.ReasonRuntimeError
	DebugReasonCompleted    = debugger.ReasonCompleted
	DebugReasonTerminated   = debugger.ReasonTerminated

	DebugBreakpointBindNextExecutableInSource   = debugger.BreakpointBindNextExecutableInSource
	DebugBreakpointBindExact                    = debugger.BreakpointBindExact
	DebugBreakpointBindNextExecutableInFunction = debugger.BreakpointBindNextExecutableInFunction
)
