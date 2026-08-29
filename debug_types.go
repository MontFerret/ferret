package ferret

import (
	"github.com/MontFerret/ferret/v2/pkg/debugger"
	"github.com/MontFerret/ferret/v2/pkg/source"
)

type (
	DebugSession               = debugger.Session
	DebugReason                = debugger.Reason
	DebugLocation              = source.Range
	DebugValue                 = debugger.Value
	DebugVariable              = debugger.Variable
	DebugFrame                 = debugger.Frame
	DebugBreakpoint            = debugger.Breakpoint
	DebugBreakpointID          = debugger.BreakpointID
	DebugBreakpointOptions     = debugger.BreakpointOptions
	DebugBreakpointBindingMode = debugger.BreakpointBindingMode
	DebugSourceLocation        = source.Location
	DebugValueReference        = debugger.ValueReference
	DebugEvent                 = debugger.Event
	DebugStateError            = debugger.StateError
	DebugFormatOptions         = debugger.FormatOptions
)

const (
	DebugReasonEntry        = debugger.ReasonEntry
	DebugReasonBreakpoint   = debugger.ReasonBreakpoint
	DebugReasonStep         = debugger.ReasonStep
	DebugReasonPause        = debugger.ReasonPause
	DebugReasonRuntimeError = debugger.ReasonRuntimeError
	DebugReasonCompleted    = debugger.ReasonCompleted
	DebugReasonTerminated   = debugger.ReasonTerminated

	DebugBreakpointBindNextExecutableInFile     = debugger.BreakpointBindNextExecutableInFile
	DebugBreakpointBindExact                    = debugger.BreakpointBindExact
	DebugBreakpointBindNextExecutableInFunction = debugger.BreakpointBindNextExecutableInFunction
)
