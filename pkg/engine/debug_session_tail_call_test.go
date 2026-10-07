package engine

import (
	"fmt"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/debugger"
	"github.com/MontFerret/ferret/v2/pkg/source"
)

const tailCallDebugQuery = `let items = [1, 2, 3]

func triple(value) {
    return value * 3
}

func double(value) {
    return triple(value * 2)
}

let mapped = (
    for item in items
        return double(item)
)

return mapped`

func newTailCallDebugSession(t *testing.T) *debugger.Session {
	t.Helper()

	engine, err := New(WithOptimizationLevel(OptimizationFull))
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() { _ = engine.Close() })
	plan, err := engine.CompileDebug(t.Context(), source.New("tail-call.fql", tailCallDebugQuery))
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() { _ = plan.Close() })
	session, err := plan.NewDebugSession(t.Context())
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() { _ = session.Close() })

	return session
}

func assertTailCallDebugFrames(t *testing.T, session *debugger.Session, names ...string) {
	t.Helper()

	frames, err := session.Frames(t.Context())
	if err != nil {
		t.Fatal(err)
	}

	if len(frames) != len(names) {
		t.Fatalf("frames = %+v, want %v", frames, names)
	}

	for i, name := range names {
		if frames[i].Name != name {
			t.Fatalf("frame %d = %q, want %q", i, frames[i].Name, name)
		}
	}
}

func TestDebugSessionTailCallFramesAcrossLoopHits(t *testing.T) {
	session := newTailCallDebugSession(t)
	breakpoint, err := session.SetBreakpoint(t.Context(), source.Location{SourceName: "tail-call.fql", Position: source.Position{Line: 4}})
	if err != nil || !breakpoint.Bound {
		t.Fatalf("breakpoint = %+v: %v", breakpoint, err)
	}

	if _, err := session.Start(t.Context()); err != nil {
		t.Fatal(err)
	}

	for hit := 1; hit <= 3; hit++ {
		event, err := session.Continue(t.Context())
		if err != nil {
			t.Fatal(err)
		}

		if event.Reason != debugger.ReasonBreakpoint || event.Location.Line != 4 || event.Depth != 2 || len(event.HitBreakpointIDs) != 1 || event.HitBreakpointIDs[0] != breakpoint.ID {
			t.Fatalf("hit %d: unexpected event %+v", hit, event)
		}

		assertTailCallDebugFrames(t, session, "triple", "double", "<main>")
		assertFrameValues(t, session, 0, map[string]string{"value": fmt.Sprint(hit * 2)})
		assertFrameValues(t, session, 1, map[string]string{"value": fmt.Sprint(hit)})
		assertFrameValues(t, session, 2, map[string]string{"item": fmt.Sprint(hit)})

		for frame, expression := range []string{"value", "value", "item"} {
			want := hit
			if frame == 0 {
				want *= 2
			}

			value, err := session.EvaluateFrame(t.Context(), frame, expression)
			if err != nil || value.Display != fmt.Sprint(want) {
				t.Fatalf("hit %d frame %d evaluation = %+v: %v", hit, frame, value, err)
			}
		}
	}

	event, err := session.Continue(t.Context())
	if err != nil || event.Reason != debugger.ReasonCompleted || event.Output == nil || string(event.Output.Data) != "[6,12,18]" {
		t.Fatalf("completion = %+v: %v", event, err)
	}
}

func pauseTailCallMain(t *testing.T, session *debugger.Session) {
	t.Helper()

	breakpoint, err := session.SetBreakpoint(t.Context(), source.Location{Position: source.Position{Line: 13}})
	if err != nil || !breakpoint.Bound {
		t.Fatalf("main breakpoint = %+v: %v", breakpoint, err)
	}

	if _, err := session.Start(t.Context()); err != nil {
		t.Fatal(err)
	}

	event, err := session.Continue(t.Context())
	if err != nil || event.Reason != debugger.ReasonBreakpoint || event.Location.Line != 13 {
		t.Fatalf("main stop = %+v: %v", event, err)
	}

	if err := session.DeleteBreakpoint(t.Context(), breakpoint.ID); err != nil {
		t.Fatal(err)
	}
}

func TestDebugSessionTailCallStepping(t *testing.T) {
	t.Run("step in and out", func(t *testing.T) {
		session := newTailCallDebugSession(t)
		pauseTailCallMain(t, session)
		for i, line := range []int{8, 4} {
			event, err := session.StepIn(t.Context())
			if err != nil || event.Reason != debugger.ReasonStep || event.Location.Line != line || event.Depth != i+1 {
				t.Fatalf("step in = %+v: %v", event, err)
			}
		}

		assertTailCallDebugFrames(t, session, "triple", "double", "<main>")

		// Neither returning UDF has another source point. The next valid caller
		// observation is the main loop's second iteration, before double executes.
		event, err := session.StepOut(t.Context())
		if err != nil || event.Reason != debugger.ReasonStep || event.Location.Line != 13 || event.Depth != 0 {
			t.Fatalf("step out = %+v: %v", event, err)
		}

		assertTailCallDebugFrames(t, session, "<main>")
		assertFrameValues(t, session, 0, map[string]string{"item": "2"})
	})

	t.Run("step over without callee breakpoint", func(t *testing.T) {
		session := newTailCallDebugSession(t)
		pauseTailCallMain(t, session)
		event, err := session.StepOver(t.Context())
		if err != nil || event.Reason != debugger.ReasonStep || event.Location.Line != 13 || event.Depth != 0 {
			t.Fatalf("step over = %+v: %v", event, err)
		}

		assertTailCallDebugFrames(t, session, "<main>")
		assertFrameValues(t, session, 0, map[string]string{"item": "2"})
	})

	t.Run("breakpoint interrupts step over", func(t *testing.T) {
		session := newTailCallDebugSession(t)
		pauseTailCallMain(t, session)
		breakpoint, err := session.SetBreakpoint(t.Context(), source.Location{Position: source.Position{Line: 4}})
		if err != nil || !breakpoint.Bound {
			t.Fatalf("callee breakpoint = %+v: %v", breakpoint, err)
		}

		event, err := session.StepOver(t.Context())
		if err != nil || event.Reason != debugger.ReasonBreakpoint || event.Location.Line != 4 || len(event.HitBreakpointIDs) != 1 || event.HitBreakpointIDs[0] != breakpoint.ID {
			t.Fatalf("interrupted step over = %+v: %v", event, err)
		}

		assertTailCallDebugFrames(t, session, "triple", "double", "<main>")
	})
}
