package ferret

import "github.com/MontFerret/ferret/v2/pkg/engine"

type (
	// DebugSession retains execution state for source-level debugging.
	DebugSession = engine.DebugSession

	// DebugReason identifies why a debug session stopped.
	DebugReason = engine.DebugReason

	// DebugLocation identifies the source range of a debug stop.
	DebugLocation = engine.DebugLocation

	// DebugValue describes a value inspected in a paused debug session.
	DebugValue = engine.DebugValue

	// DebugVariable describes a visible variable in a paused debug session.
	DebugVariable = engine.DebugVariable

	// DebugFrame describes a retained call frame.
	DebugFrame = engine.DebugFrame

	// DebugBreakpoint describes a requested breakpoint and its executable binding.
	DebugBreakpoint = engine.DebugBreakpoint

	// DebugBreakpointID identifies a breakpoint within its debug session.
	DebugBreakpointID = engine.DebugBreakpointID

	// DebugBreakpointOptions configures how a requested source location binds.
	DebugBreakpointOptions = engine.DebugBreakpointOptions

	// DebugBreakpointBindingMode selects how a requested source location is bound.
	DebugBreakpointBindingMode = engine.DebugBreakpointBindingMode

	// DebugSourceLocation identifies a position requested for a breakpoint.
	DebugSourceLocation = engine.DebugSourceLocation

	// DebugValueReference identifies an expandable value in the current paused state.
	DebugValueReference = engine.DebugValueReference

	// DebugEvent describes a debug stop, failure, or completed output.
	DebugEvent = engine.DebugEvent

	// DebugStateError reports an operation that is invalid in the session's current state.
	DebugStateError = engine.DebugStateError

	// DebugFormatOptions bounds debugger value formatting.
	DebugFormatOptions = engine.DebugFormatOptions
)

const (
	DebugReasonEntry        = engine.DebugReasonEntry
	DebugReasonBreakpoint   = engine.DebugReasonBreakpoint
	DebugReasonStep         = engine.DebugReasonStep
	DebugReasonPause        = engine.DebugReasonPause
	DebugReasonRuntimeError = engine.DebugReasonRuntimeError
	DebugReasonCompleted    = engine.DebugReasonCompleted
	DebugReasonTerminated   = engine.DebugReasonTerminated

	DebugBreakpointBindNextExecutableInSource   = engine.DebugBreakpointBindNextExecutableInSource
	DebugBreakpointBindExact                    = engine.DebugBreakpointBindExact
	DebugBreakpointBindNextExecutableInFunction = engine.DebugBreakpointBindNextExecutableInFunction
)
