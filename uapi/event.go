package uapi

import apidebugger "github.com/MontFerret/api/debugger"

func convertEvent(event *apidebugger.Event) *apidebugger.Event {
	if event != nil {
		// Native returns a fresh event, including its mutable fields.
		event.Error = wrapDiagnosticError(event.Error)
	}

	return event
}
