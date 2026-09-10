package uapi

import (
	"github.com/MontFerret/api"
	apidebugger "github.com/MontFerret/api/debugger"
)

func outputValue(output *api.Output) api.Output {
	if output == nil {
		return api.Output{}
	}

	// Native has already transferred the encoded bytes to the caller.
	return *output
}

func convertEvent(event *apidebugger.Event) *apidebugger.Event {
	if event != nil {
		// Native returns a fresh event, including its mutable fields.
		event.Error = wrapDiagnosticError(event.Error)
	}

	return event
}
