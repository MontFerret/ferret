package uapi

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/MontFerret/api"
	apidebugger "github.com/MontFerret/api/debugger"
	apidiagnostics "github.com/MontFerret/api/diagnostics"
	nativedebug "github.com/MontFerret/ferret/v2/pkg/debugger"
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

	breakpoint, err := session.SetBreakpointAt(context.Background(), location, apidebugger.BreakpointOptions{BindingMode: apidebugger.BreakpointBindExact})
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

			breakpoint, err := session.SetBreakpoint(context.Background(), api.Location{SourceName: "buffer://steps", Position: api.Position{Line: line, Column: 1}})
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

				if value.IndexedCount == nil || *value.IndexedCount != 2 || value.NamedCount == nil || *value.NamedCount != 0 {
					t.Fatalf("array counts: %+v", value)
				}

				if page, err := session.Variables(t.Context(), reference, apidebugger.VariablePage{Count: -1}); page != nil || !errors.Is(err, apidebugger.ErrInvalidVariablePage) {
					t.Fatalf("invalid page: %v, %v", page, err)
				}

				children, err := session.Variables(t.Context(), reference, apidebugger.VariablePage{Filter: apidebugger.VariableFilterIndexed, Start: 1, Count: 1})
				if err != nil || len(children) != 1 || children[0].Name != "1" || children[0].Value.Display != "2" {
					t.Fatalf("children=%+v err=%v", children, err)
				}
			}

			if err := session.DeleteBreakpoint(context.Background(), breakpoint.ID); err != nil {
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
				if _, err := session.Variables(context.Background(), reference, apidebugger.VariablePage{}); !errors.Is(err, apidebugger.ErrInvalidValueReference) {
					t.Fatal("resumed session accepted a stale value reference")
				}
			}

			completed, err := session.Continue(t.Context())
			if err != nil || completed.Reason != apidebugger.ReasonCompleted || completed.Output == nil || string(completed.Output.Data) != "[[1,2],3]" {
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
			runtime := Wrap(engine, "test-core-version")
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
				if err := session.Pause(context.Background()); err != nil {
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

	plan, err := runtime.CompileDebug(t.Context(), api.NewSource("buffer://error", "LET box = [1]\nLET value = box[0]\nRETURN value / @zero"))
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

	if frames, err := session.Frames(context.Background()); err != nil || len(frames) == 0 {
		t.Fatalf("frames=%+v err=%v", frames, err)
	}

	if value, err := session.EvaluateFrame(t.Context(), 0, "value + 1"); err != nil || value.Display != "2" {
		t.Fatalf("evaluation at runtime error stop: value=%+v err=%v", value, err)
	}

	locals, err := session.Locals(t.Context())
	if err != nil {
		t.Fatal(err)
	}

	var ref apidebugger.ValueReference

	for _, binding := range locals {
		if binding.Name == "box" {
			ref = binding.Value.Reference
		}
	}

	for range 2 {
		page, err := session.Variables(t.Context(), ref, apidebugger.VariablePage{})
		if err != nil || len(page) != 1 || page[0].Value.Display != "1" {
			t.Fatalf("runtime-error page: %+v, %v", page, err)
		}
	}

	if event, err := session.Continue(t.Context()); err != nil || event.Reason != apidebugger.ReasonTerminated {
		t.Fatalf("terminate: %+v, %v", event, err)
	}

	if page, err := session.Variables(t.Context(), ref, apidebugger.VariablePage{}); page != nil || !errors.Is(err, &nativedebug.StateError{}) {
		t.Fatalf("terminal state: %+v, %v", page, err)
	}
}

func TestDebugSummaryOnlyObjectReference(t *testing.T) {
	rt := newTestRuntime(t)
	plan, err := rt.CompileDebug(t.Context(), api.NewAnonymousSource("RETURN @object"))
	if err != nil {
		t.Fatal(err)
	}

	defer plan.Close()
	object := runtime.NewObjectOf(256)
	object.SetString("only", runtime.True)
	session, err := plan.NewDebugSession(t.Context(), api.WithParam("object", object))
	if err != nil {
		t.Fatal(err)
	}

	defer session.Close()
	if _, err := session.Start(t.Context()); err != nil {
		t.Fatal(err)
	}

	value, err := session.Evaluate(t.Context(), "@object")
	if err != nil || value.Display != "Object(1)" || value.Reference != 0 ||
		value.NamedCount == nil || *value.NamedCount != 1 || value.IndexedCount == nil || *value.IndexedCount != 0 {
		t.Fatalf("summary: %+v, %v", value, err)
	}

	page, err := session.Variables(t.Context(), value.Reference, apidebugger.VariablePage{Filter: apidebugger.VariableFilterIndexed})
	if page != nil || !errors.Is(err, apidebugger.ErrInvalidValueReference) {
		t.Fatalf("zero reference: %v, %v", page, err)
	}
}

type filteredDebugValue struct{ runtime.Value }

func (filteredDebugValue) DebugChildrenInfo() runtime.DebugChildrenInfo {
	return runtime.DebugChildrenInfo{Indexed: runtime.DebugChildView{Readable: true}}
}

func (filteredDebugValue) DebugChildAt(_ context.Context, kind apidebugger.VariableFilter, position int64) (string, runtime.Value, bool, error) {
	if kind != apidebugger.VariableFilterIndexed {
		return "", nil, false, errors.ErrUnsupported
	}

	return "", runtime.True, position == 0, nil
}

func TestDebugFilteredOnlyProvider(t *testing.T) {
	rt := newTestRuntime(t)
	plan, err := rt.CompileDebug(t.Context(), api.NewAnonymousSource("RETURN @view"))
	if err != nil {
		t.Fatal(err)
	}

	defer plan.Close()
	session, err := plan.NewDebugSession(t.Context(), api.WithParam("view", filteredDebugValue{}))
	if err != nil {
		t.Fatal(err)
	}

	defer session.Close()
	if _, err := session.Start(t.Context()); err != nil {
		t.Fatal(err)
	}

	value, err := session.Evaluate(t.Context(), "@view")
	if err != nil || !value.Reference.Valid() || value.NamedCount != nil || value.IndexedCount != nil {
		t.Fatalf("view: %+v %v", value, err)
	}

	if page, err := session.Variables(t.Context(), value.Reference, apidebugger.VariablePage{}); page != nil || !errors.Is(err, errors.ErrUnsupported) {
		t.Fatalf("combined: %+v %v", page, err)
	}

	page, err := session.Variables(t.Context(), value.Reference, apidebugger.VariablePage{Filter: apidebugger.VariableFilterIndexed})
	if err != nil || len(page) != 1 || page[0].Name != "0" || page[0].Value.Display != "true" {
		t.Fatalf("indexed: %+v, %v", page, err)
	}
}
