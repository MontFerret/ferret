package uapi

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/MontFerret/api"
	apidebugger "github.com/MontFerret/api/debugger"
	apisource "github.com/MontFerret/api/source"
)

func TestDebugFramePositionsAddressLocalsAndEvaluation(t *testing.T) {
	runtime := newTestRuntime(t)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	t.Cleanup(cancel)

	sourcePath := filepath.Join(t.TempDir(), "query.fql")
	// Keep calls out of tail position so all three frames remain inspectable.
	plan, err := runtime.CompileDebug(ctx, api.NewSource(sourcePath, `LET marker = 10
FUNC outer(marker) {
  FUNC inner(marker) {
    RETURN marker
  }

  LET result = inner(30)
  RETURN result + marker
}

RETURN outer(20) + marker`))
	if err != nil {
		t.Fatalf("CompileDebug: %v", err)
	}

	t.Cleanup(func() {
		if err := plan.Close(); err != nil {
			t.Errorf("plan Close: %v", err)
		}
	})

	session, err := plan.NewDebugSession(ctx)
	if err != nil {
		t.Fatalf("NewDebugSession: %v", err)
	}

	t.Cleanup(func() {
		if err := session.Close(); err != nil {
			t.Errorf("session Close: %v", err)
		}
	})

	breakpoint, err := session.SetBreakpoint(apisource.Location{
		SourceName: sourcePath,
		Position:   apisource.Position{Line: 4, Column: 1},
	})
	if err != nil {
		t.Fatalf("SetBreakpoint: %v", err)
	}

	if !breakpoint.Bound || breakpoint.Location.Line != 4 {
		t.Fatalf("breakpoint = %+v, want bound at line 4", breakpoint)
	}

	entry, err := session.Start(ctx)
	if err != nil {
		t.Fatalf("Start: %v", err)
	}

	if entry == nil || entry.Reason != apidebugger.ReasonEntry {
		t.Fatalf("entry event = %+v", entry)
	}

	stopped, err := session.Continue(ctx)
	if err != nil {
		t.Fatalf("Continue: %v", err)
	}

	if stopped == nil || stopped.Reason != apidebugger.ReasonBreakpoint || stopped.Location.Line != 4 {
		t.Fatalf("breakpoint event = %+v", stopped)
	}

	frames, err := session.Frames()
	if err != nil {
		t.Fatalf("Frames: %v", err)
	}

	wantNames := []string{"inner", "outer", "<main>"}
	wantDisplays := []string{"30", "20", "10"}

	if len(frames) != len(wantNames) {
		t.Fatalf("Frames = %+v, want %d frames", frames, len(wantNames))
	}

	for index, frame := range frames {
		if frame.Name != wantNames[index] || frame.Location.SourceName != sourcePath {
			t.Fatalf("Frames()[%d] = %+v, want %s in %s", index, frame, wantNames[index], sourcePath)
		}

		locals, err := session.FrameLocals(index)
		if err != nil {
			t.Fatalf("FrameLocals(%d): %v", index, err)
		}

		var marker apidebugger.Value
		for _, local := range locals {
			if local.Name == "marker" && !local.Param {
				marker = local.Value

				break
			}
		}

		if marker.Display != wantDisplays[index] {
			t.Fatalf("FrameLocals(%d) = %+v, want marker = %s", index, locals, wantDisplays[index])
		}

		value, err := session.EvaluateFrame(ctx, index, "marker")
		if err != nil {
			t.Fatalf("EvaluateFrame(%d): %v", index, err)
		}

		if value != marker {
			t.Fatalf("EvaluateFrame(%d) = %+v, want local value %+v", index, value, marker)
		}
	}
}
