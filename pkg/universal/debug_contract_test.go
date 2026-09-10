package universal

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/MontFerret/api"
	apidebugger "github.com/MontFerret/api/debugger"
	apidiagnostics "github.com/MontFerret/api/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/engine"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func TestDebugLocationsUseByteCoordinates(t *testing.T) {
	runtime := newTestRuntime(t)
	const query = `LET prefix = "é" RETURN prefix`
	src := api.NewSource("cell:unicode", query)

	plan, err := runtime.CompileDebug(t.Context(), src)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() { _ = plan.Close() })

	session, err := plan.NewDebugSession(t.Context())
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() { _ = session.Close() })
	start := strings.Index(query, "RETURN")
	location := api.Location{SourceName: src.Name, Position: api.Position{Line: 1, Column: start + 1}}

	breakpoint, err := session.SetBreakpointAt(location, apidebugger.BreakpointOptions{BindingMode: apidebugger.BreakpointBindExact})
	if err != nil || !breakpoint.Bound {
		t.Fatalf("breakpoint=%+v error=%v", breakpoint, err)
	}

	if _, err := session.Start(t.Context()); err != nil {
		t.Fatal(err)
	}

	event, err := session.Continue(t.Context())
	if err != nil || event.Reason != apidebugger.ReasonBreakpoint {
		t.Fatalf("event=%+v error=%v", event, err)
	}

	want := api.Range{Location: location, Span: api.Span{Start: start, End: len(query)}}
	if event.Location != want {
		t.Fatalf("debug location=%+v, want %+v", event.Location, want)
	}
}

func TestDebugCommandsAndReferenceLifetime(t *testing.T) {
	runtime := newTestRuntime(t)

	plan, err := runtime.CompileDebug(t.Context(), api.NewSource("buffer://steps", `LET value = [1, 2]
FUNC inner(input) {
  RETURN input + 1
}

LET result = inner(2)
RETURN [value, result]`))
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() { _ = plan.Close() })
	for _, step := range []string{"in", "over", "out"} {
		t.Run(step, func(t *testing.T) {
			session, err := plan.NewDebugSession(t.Context())
			if err != nil {
				t.Fatal(err)
			}

			t.Cleanup(func() { _ = session.Close() })
			line := 5

			if step == "out" {
				line = 3
			}

			breakpoint, err := session.SetBreakpoint(api.Location{SourceName: "buffer://steps", Position: api.Position{Line: line, Column: 1}})
			if err != nil || !breakpoint.Bound {
				t.Fatalf("breakpoint=%+v err=%v", breakpoint, err)
			}

			if _, err := session.Start(t.Context()); err != nil {
				t.Fatal(err)
			}

			stop, err := session.Continue(t.Context())
			if err != nil || stop.Reason != apidebugger.ReasonBreakpoint || len(stop.HitBreakpointIDs) != 1 || stop.HitBreakpointIDs[0] != breakpoint.ID {
				t.Fatalf("stop=%+v err=%v", stop, err)
			}

			var reference apidebugger.ValueReference

			if step != "out" {
				value, err := session.Evaluate(t.Context(), "value")
				if err != nil || !value.Reference.Valid() {
					t.Fatalf("value=%+v err=%v", value, err)
				}

				reference = value.Reference

				children, err := session.Variables(reference)
				if err != nil || len(children) != 2 {
					t.Fatalf("children=%+v err=%v", children, err)
				}
			}

			if err := session.DeleteBreakpoint(breakpoint.ID); err != nil {
				t.Fatal(err)
			}

			command := session.StepIn

			switch step {
			case "over":
				command = session.StepOver
			case "out":
				command = session.StepOut
			}

			event, err := command(t.Context())
			if err != nil || event.Reason != apidebugger.ReasonStep {
				t.Fatalf("step=%+v err=%v", event, err)
			}

			if step == "in" && event.Depth <= stop.Depth || step == "over" && event.Depth > stop.Depth || step == "out" && event.Depth >= stop.Depth {
				t.Fatalf("%s depth %d -> %d", step, stop.Depth, event.Depth)
			}

			if reference.Valid() {
				if _, err := session.Variables(reference); err == nil {
					t.Fatal("resumed session accepted a stale value reference")
				}
			}

			completed, err := session.Continue(t.Context())
			if err != nil || completed.Reason != apidebugger.ReasonCompleted || completed.Output == nil || string(completed.Output.Content) != "[[1,2],3]" {
				t.Fatalf("completion=%+v err=%v", completed, err)
			}
		})
	}
}

func TestDebugPauseAndCancellationReachNativeExecution(t *testing.T) {
	for _, mode := range []string{"pause", "cancel", "close"} {
		t.Run(mode, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
			t.Cleanup(cancel)
			entered, release := make(chan struct{}), make(chan struct{})

			engine, err := engine.New(engine.WithFunctionsRegistrar(func(ns runtime.Namespace) {
				ns.Function().A0().Add("WAIT_NATIVE", func(ctx context.Context) (runtime.Value, error) {
					close(entered)
					select {
					case <-release:
						return runtime.NewInt(1), nil
					case <-ctx.Done():
						return nil, ctx.Err()
					}
				})
			}))
			if err != nil {
				t.Fatal(err)
			}

			t.Cleanup(func() { _ = engine.Close() })
			runtime := Wrap(engine)
			t.Cleanup(func() { _ = runtime.Close() })

			plan, err := runtime.CompileDebug(ctx, api.NewAnonymousSource("LET value = WAIT_NATIVE()\nRETURN value"))
			if err != nil {
				t.Fatal(err)
			}

			t.Cleanup(func() { _ = plan.Close() })

			session, err := plan.NewDebugSession(ctx)
			if err != nil {
				t.Fatal(err)
			}

			t.Cleanup(func() { _ = session.Close() })

			if _, err := session.Start(ctx); err != nil {
				t.Fatal(err)
			}

			type response struct {
				event *apidebugger.Event
				err   error
			}

			result := make(chan response, 1)
			commandCtx, commandCancel := context.WithCancel(ctx)
			t.Cleanup(commandCancel)
			go func() { event, err := session.Continue(commandCtx); result <- response{event, err} }()
			select {
			case <-entered:
			case <-ctx.Done():
				t.Fatal(ctx.Err())
			}

			switch mode {
			case "pause":
				if err := session.Pause(); err != nil {
					t.Fatal(err)
				}

				close(release)
			case "cancel":
				commandCancel()
			case "close":
				if err := session.Close(); err != nil {
					t.Fatal(err)
				}
			}

			select {
			case response := <-result:
				if mode == "pause" {
					if response.err != nil || response.event == nil || response.event.Reason != apidebugger.ReasonPause {
						t.Fatalf("pause=%+v", response)
					}
				} else if response.event == nil || response.event.Reason != apidebugger.ReasonTerminated || !errors.Is(errors.Join(response.err, response.event.Error), context.Canceled) {
					t.Fatalf("termination event=%+v err=%v", response.event, response.err)
				}
			case <-ctx.Done():
				t.Fatal(ctx.Err())
			}
		})
	}
}

func TestDebugRuntimeErrorRemainsInspectable(t *testing.T) {
	runtime := newTestRuntime(t)

	plan, err := runtime.CompileDebug(t.Context(), api.NewSource("buffer://error", "LET value = 1\nRETURN value / @zero"))
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() { _ = plan.Close() })

	session, err := plan.NewDebugSession(t.Context(), api.WithParam("zero", 0))
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() { _ = session.Close() })

	if _, err := session.Start(t.Context()); err != nil {
		t.Fatal(err)
	}

	event, err := session.Continue(t.Context())

	var diagnostics apidiagnostics.Diagnostics
	if err != nil || event.Reason != apidebugger.ReasonRuntimeError || !errors.As(event.Error, &diagnostics) || len(diagnostics) != 1 || diagnostics[0].Source.Name != "buffer://error" {
		t.Fatalf("event=%+v err=%v diagnostics=%+v", event, err, diagnostics)
	}

	if frames, err := session.Frames(); err != nil || len(frames) == 0 {
		t.Fatalf("frames=%+v err=%v", frames, err)
	}
}
