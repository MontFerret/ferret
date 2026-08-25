package ferret_test

import (
	"testing"

	"github.com/MontFerret/api"
	apidebugger "github.com/MontFerret/api/debugger"
	apisource "github.com/MontFerret/api/source"
	ferret "github.com/MontFerret/ferret/v2"
)

func TestUniversalDebugSessionBridge(t *testing.T) {
	t.Parallel()

	engine, err := ferret.New()
	if err != nil {
		t.Fatalf("new engine: %v", err)
	}
	t.Cleanup(func() { _ = engine.Close() })

	var universal api.Runtime = engine
	plan, err := universal.CompileDebug(
		t.Context(),
		apisource.File{
			Name:    "bridge.fql",
			Content: "LET value = 1\nRETURN value + 1",
		},
		api.WithOptimizationLevel(api.OptimizationAggressive),
	)
	if err != nil {
		t.Fatalf("compile debug: %v", err)
	}
	t.Cleanup(func() { _ = plan.Close() })

	session, err := plan.NewDebugSession(t.Context())
	if err != nil {
		t.Fatalf("new debug session: %v", err)
	}
	t.Cleanup(func() { _ = session.Close() })

	breakpoint, err := session.SetBreakpoint(apisource.Location{
		File: "bridge.fql",
		Position: apisource.Position{
			Line: 2,
		},
	})
	if err != nil {
		t.Fatalf("set breakpoint: %v", err)
	}

	if !breakpoint.Bound || breakpoint.Location.Line != 2 {
		t.Fatalf("breakpoint = %#v, want bound line 2", breakpoint)
	}

	event, err := session.Start(t.Context())
	if err != nil {
		t.Fatalf("start: %v", err)
	}

	if event.Reason != apidebugger.ReasonEntry || event.Location.File != "bridge.fql" {
		t.Fatalf("entry event = %#v", event)
	}

	event, err = session.Continue(t.Context())
	if err != nil {
		t.Fatalf("continue to breakpoint: %v", err)
	}

	if event.Reason != apidebugger.ReasonBreakpoint || event.Location.Line != 2 {
		t.Fatalf("breakpoint event = %#v", event)
	}

	event, err = session.Continue(t.Context())
	if err != nil {
		t.Fatalf("continue to completion: %v", err)
	}

	if event.Reason != apidebugger.ReasonCompleted || event.Output == nil || string(event.Output.Content) != "2" {
		t.Fatalf("completion event = %#v", event)
	}
}
