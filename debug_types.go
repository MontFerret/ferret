package ferret

import "github.com/MontFerret/ferret/v2/pkg/debugger"

type (
	DebugSession       = debugger.Session
	DebugReason        = debugger.Reason
	DebugLocation      = debugger.Location
	DebugValue         = debugger.Value
	DebugVariable      = debugger.Variable
	DebugFrame         = debugger.Frame
	DebugBreakpoint    = debugger.Breakpoint
	DebugEvent         = debugger.Event
	DebugStateError    = debugger.StateError
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
)
