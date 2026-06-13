package ferret

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/source"
)

func TestDebugSessionBreakpointsLocalsEvaluateAndComplete(t *testing.T) {
	engine, err := New()
	if err != nil {
		t.Fatal(err)
	}
	defer engine.Close()

	src := source.New("debug.fql", "LET x = 1\n\nVAR y = 2\ny = y + x\nRETURN y")
	plan, err := engine.CompileDebug(context.Background(), src)
	if err != nil {
		t.Fatal(err)
	}
	defer plan.Close()

	session, err := plan.NewDebugSession(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	defer session.Close()

	breakpoint, err := session.SetBreakpoint("debug.fql", 2)
	if err != nil {
		t.Fatal(err)
	}
	if !breakpoint.Bound || breakpoint.Line != 3 {
		t.Fatalf("unexpected breakpoint: %#v", breakpoint)
	}

	event, err := session.Start(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if event.Reason != DebugReasonEntry || event.Location.Line != 1 {
		t.Fatalf("unexpected entry event: %#v", event)
	}
	locals, err := session.Locals()
	if err != nil {
		t.Fatal(err)
	}
	if len(locals) != 0 {
		t.Fatalf("declaration must not be visible before execution: %#v", locals)
	}

	event, err = session.Continue(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if event.Reason != DebugReasonBreakpoint || event.Location.Line != 3 {
		t.Fatalf("unexpected breakpoint event: %#v", event)
	}
	locals, err = session.Locals()
	if err != nil {
		t.Fatal(err)
	}
	if len(locals) != 1 || locals[0].Name != "x" || locals[0].Value.Display != "1" {
		t.Fatalf("unexpected locals: %#v", locals)
	}

	event, err = session.Step(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if event.Location.Line != 4 {
		t.Fatalf("unexpected step location: %#v", event)
	}
	value, err := session.Evaluate(context.Background(), "x + y")
	if err != nil {
		t.Fatal(err)
	}
	if value.Display != "3" {
		t.Fatalf("unexpected evaluated value: %#v", value)
	}
	if _, err := session.Evaluate(context.Background(), "LENGTH([1])"); err == nil {
		t.Fatal("expected function call evaluation to be rejected")
	}

	event, err = session.Continue(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if event.Reason != DebugReasonCompleted || event.Output == nil || string(event.Output.Content) != "3" {
		t.Fatalf("unexpected completion event: %#v", event)
	}
	if _, err := session.Continue(context.Background()); err == nil || !errors.Is(err, &DebugStateError{}) {
		t.Fatalf("expected typed invalid-state error, got %v", err)
	}
}

func TestDebugSessionBreakpointBindsOnePointPerLine(t *testing.T) {
	engine, err := New()
	if err != nil {
		t.Fatal(err)
	}
	defer engine.Close()

	plan, err := engine.CompileDebug(context.Background(), source.New("same-line.fql", "LET x = 1 RETURN x"))
	if err != nil {
		t.Fatal(err)
	}
	defer plan.Close()
	session, err := plan.NewDebugSession(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	defer session.Close()

	breakpoint, err := session.SetBreakpoint("same-line.fql", 1)
	if err != nil {
		t.Fatal(err)
	}
	if !breakpoint.Bound {
		t.Fatalf("expected bound breakpoint: %#v", breakpoint)
	}
	if breakpoint.RequestedColumn != 0 || breakpoint.FunctionID != -1 {
		t.Fatalf("unexpected same-line breakpoint identity: %#v", breakpoint)
	}
	if _, err := session.Start(context.Background()); err != nil {
		t.Fatal(err)
	}
	event, err := session.Continue(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if event.Reason != DebugReasonCompleted {
		t.Fatalf("expected breakpoint to bind only the first same-line point, got %#v", event)
	}
}

func TestDebugSessionBreakpointLifecyclePreservesIDs(t *testing.T) {
	engine, err := New()
	if err != nil {
		t.Fatal(err)
	}
	defer engine.Close()

	plan, err := engine.CompileDebug(context.Background(), source.New("breakpoints.fql", "RETURN 1"))
	if err != nil {
		t.Fatal(err)
	}
	defer plan.Close()
	session, err := plan.NewDebugSession(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	defer session.Close()

	first, err := session.SetBreakpoint("other.fql", 1)
	if err != nil {
		t.Fatal(err)
	}
	if first.Bound {
		t.Fatalf("expected unbound breakpoint: %#v", first)
	}
	if err := session.DeleteBreakpoint(first.ID); err != nil {
		t.Fatal(err)
	}
	second, err := session.SetBreakpoint("breakpoints.fql", 1)
	if err != nil {
		t.Fatal(err)
	}
	if !second.Bound || second.ID <= first.ID {
		t.Fatalf("unexpected replacement breakpoint: %#v after %#v", second, first)
	}
	if got := session.Breakpoints(); len(got) != 1 || got[0].ID != second.ID {
		t.Fatalf("unexpected breakpoint snapshot: %#v", got)
	}
}

func TestDebugSessionLocalsIncludeDeclaredBindParameters(t *testing.T) {
	engine, err := New()
	if err != nil {
		t.Fatal(err)
	}
	defer engine.Close()

	plan, err := engine.CompileDebug(context.Background(), source.New("params.fql", "RETURN @input"))
	if err != nil {
		t.Fatal(err)
	}
	defer plan.Close()
	session, err := plan.NewDebugSession(
		context.Background(),
		WithSessionParam("input", 4),
		WithSessionParam("unrelated", 9),
	)
	if err != nil {
		t.Fatal(err)
	}
	defer session.Close()

	if _, err := session.Start(context.Background()); err != nil {
		t.Fatal(err)
	}
	locals, err := session.Locals()
	if err != nil {
		t.Fatal(err)
	}
	if len(locals) != 1 || locals[0].Name != "@input" || !locals[0].Param || locals[0].Value.Display != "4" {
		t.Fatalf("unexpected bind parameters: %#v", locals)
	}
	value, err := session.Evaluate(context.Background(), "@input + 1")
	if err != nil || value.Display != "5" {
		t.Fatalf("unexpected parameter evaluation: %#v, %v", value, err)
	}
}

func TestDebugPlanRunsNormally(t *testing.T) {
	engine, err := New()
	if err != nil {
		t.Fatal(err)
	}
	defer engine.Close()

	plan, err := engine.CompileDebug(context.Background(), source.New("normal.fql", "LET x = 2\nRETURN x + 1"))
	if err != nil {
		t.Fatal(err)
	}
	defer plan.Close()
	session, err := plan.NewSession(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	defer session.Close()

	output, err := session.Run(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if string(output.Content) != "3" {
		t.Fatalf("unexpected output: %s", output.Content)
	}
}

func TestDebugSessionPauseAndSafeObjectInspection(t *testing.T) {
	engine, err := New()
	if err != nil {
		t.Fatal(err)
	}
	defer engine.Close()
	plan, err := engine.CompileDebug(context.Background(), source.New("object.fql", "LET obj = {b: 2, a: 1}\nRETURN obj"))
	if err != nil {
		t.Fatal(err)
	}
	defer plan.Close()
	session, err := plan.NewDebugSession(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	defer session.Close()

	if _, err := session.Start(context.Background()); err != nil {
		t.Fatal(err)
	}
	if err := session.Pause(); err != nil {
		t.Fatal(err)
	}
	event, err := session.Continue(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if event.Reason != DebugReasonPause || event.Location.Line != 2 {
		t.Fatalf("unexpected pause event: %#v", event)
	}
	locals, err := session.Locals()
	if err != nil {
		t.Fatal(err)
	}
	if len(locals) != 1 || locals[0].Value.Display != `{"a": 1, "b": 2}` {
		t.Fatalf("unexpected formatted object: %#v", locals)
	}
	value, err := session.Evaluate(context.Background(), "obj.a")
	if err != nil || value.Display != "1" {
		t.Fatalf("unexpected member evaluation: %#v, %v", value, err)
	}
}

func TestDebugSessionRuntimeErrorPreservesLocals(t *testing.T) {
	engine, err := New()
	if err != nil {
		t.Fatal(err)
	}
	defer engine.Close()

	plan, err := engine.CompileDebug(context.Background(), source.New("error.fql", "LET x = 7\nRETURN x / 0"))
	if err != nil {
		t.Fatal(err)
	}
	defer plan.Close()
	session, err := plan.NewDebugSession(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	defer session.Close()

	if _, err := session.Start(context.Background()); err != nil {
		t.Fatal(err)
	}
	event, err := session.Continue(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if event.Reason != DebugReasonRuntimeError || event.Error == nil {
		t.Fatalf("expected runtime error pause, got %#v", event)
	}
	if _, ok := event.Error.(diagnostics.Formattable); !ok {
		t.Fatalf("expected formattable runtime error, got %T", event.Error)
	}
	formatted := diagnostics.Format(event.Error)
	if !strings.Contains(formatted, "error.fql:2") || !strings.Contains(formatted, "RETURN x / 0") {
		t.Fatalf("expected formatted source diagnostic, got:\n%s", formatted)
	}
	locals, err := session.Locals()
	if err != nil {
		t.Fatal(err)
	}
	if len(locals) != 1 || locals[0].Name != "x" || locals[0].Value.Display != "7" {
		t.Fatalf("unexpected post-error locals: %#v", locals)
	}
	event, err = session.Continue(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if event.Reason != DebugReasonTerminated || event.Error == nil || !strings.Contains(event.Error.Error(), "division") {
		t.Fatalf("unexpected termination: %#v", event)
	}
	if _, ok := event.Error.(diagnostics.Formattable); !ok {
		t.Fatalf("expected formattable termination error, got %T", event.Error)
	}
}

func TestDebugSessionRuntimeErrorJoinsAfterRunHookFailure(t *testing.T) {
	afterErr := errors.New("after run failed")
	engine, err := New(WithAfterRunHook(func(context.Context, error) error {
		return afterErr
	}))
	if err != nil {
		t.Fatal(err)
	}
	defer engine.Close()

	plan, err := engine.CompileDebug(context.Background(), source.New("error.fql", "RETURN 1 / 0"))
	if err != nil {
		t.Fatal(err)
	}
	defer plan.Close()
	session, err := plan.NewDebugSession(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	defer session.Close()

	if _, err := session.Start(context.Background()); err != nil {
		t.Fatal(err)
	}
	event, err := session.Continue(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if event.Reason != DebugReasonRuntimeError || !errors.Is(event.Error, afterErr) {
		t.Fatalf("expected joined after-run hook failure, got %#v", event)
	}
	var runtimeErr diagnostics.FormattableError
	if !errors.As(event.Error, &runtimeErr) {
		t.Fatalf("expected joined runtime diagnostic, got %T", event.Error)
	}
}

func TestDebugSessionResumePreservesBeforeRunContextValues(t *testing.T) {
	type contextKey string

	const key contextKey = "debug-session-value"
	engine, err := New(
		WithBeforeRunHook(func(ctx context.Context) (context.Context, error) {
			return context.WithValue(ctx, key, "hook-value"), nil
		}),
		WithFunctionsRegistrar(func(fns runtime.FunctionDefs) {
			fns.A0().Add("DEBUG_CONTEXT_VALUE", func(ctx context.Context) (runtime.Value, error) {
				value, _ := ctx.Value(key).(string)
				return runtime.NewString(value), nil
			})
		}),
	)
	if err != nil {
		t.Fatal(err)
	}
	defer engine.Close()

	plan, err := engine.CompileDebug(
		context.Background(),
		source.New("context.fql", "LET x = 1\nRETURN DEBUG_CONTEXT_VALUE()"),
	)
	if err != nil {
		t.Fatal(err)
	}
	defer plan.Close()

	session, err := plan.NewDebugSession(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	defer session.Close()

	if _, err := session.Start(context.Background()); err != nil {
		t.Fatal(err)
	}
	event, err := session.Continue(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if event.Reason != DebugReasonCompleted || event.Output == nil || string(event.Output.Content) != `"hook-value"` {
		t.Fatalf("unexpected completion event: %#v", event)
	}
}

func TestPlanNewDebugSessionRequiresDebugCompilation(t *testing.T) {
	engine, err := New()
	if err != nil {
		t.Fatal(err)
	}
	defer engine.Close()
	plan, err := engine.Compile(context.Background(), source.NewAnonymous("RETURN 1"))
	if err != nil {
		t.Fatal(err)
	}
	defer plan.Close()
	if _, err := plan.NewDebugSession(context.Background()); err == nil {
		t.Fatal("expected non-debug plan to be rejected")
	}
}

func TestDebugSessionStepIntoAndOut(t *testing.T) {
	engine, err := New()
	if err != nil {
		t.Fatal(err)
	}
	defer engine.Close()
	query := `FUNC add(a) (
  LET b = a + 1
  RETURN b
)
LET x = add(2)
RETURN x`
	plan, err := engine.CompileDebug(context.Background(), source.New("udf.fql", query))
	if err != nil {
		t.Fatal(err)
	}
	defer plan.Close()
	session, err := plan.NewDebugSession(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	defer session.Close()

	event, err := session.Start(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if event.Location.Line != 5 {
		t.Fatalf("unexpected entry: %#v", event)
	}
	event, err = session.Step(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if event.Location.Line != 2 || event.Depth != 1 {
		t.Fatalf("expected step into UDF, got %#v", event)
	}
	frames, err := session.Frames()
	if err != nil {
		t.Fatal(err)
	}
	if len(frames) != 2 || frames[0].Name != "add" || frames[1].Name != "<main>" {
		t.Fatalf("unexpected frames: %#v", frames)
	}
	event, err = session.Out(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if event.Location.Line != 6 || event.Depth != 0 {
		t.Fatalf("expected step out to main, got %#v", event)
	}
}

func TestDebugSessionNextStepsOverCallAndOutFromMainCompletes(t *testing.T) {
	engine, err := New()
	if err != nil {
		t.Fatal(err)
	}
	defer engine.Close()
	query := `FUNC add(a) (
  RETURN a + 1
)
LET x = add(2)
RETURN x`
	plan, err := engine.CompileDebug(context.Background(), source.New("next.fql", query))
	if err != nil {
		t.Fatal(err)
	}
	defer plan.Close()
	session, err := plan.NewDebugSession(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	defer session.Close()

	event, err := session.Start(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if event.Location.Line != 4 {
		t.Fatalf("unexpected entry: %#v", event)
	}
	event, err = session.Next(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if event.Location.Line != 5 || event.Depth != 0 {
		t.Fatalf("expected next to step over UDF, got %#v", event)
	}
	event, err = session.Out(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if event.Reason != DebugReasonCompleted || event.Output == nil || string(event.Output.Content) != "3" {
		t.Fatalf("expected out from main to complete, got %#v", event)
	}
}

func TestDebugSessionStopsOnRepeatedLoopLocation(t *testing.T) {
	engine, err := New()
	if err != nil {
		t.Fatal(err)
	}
	defer engine.Close()
	plan, err := engine.CompileDebug(context.Background(), source.New("loop.fql", "FOR i IN 1..2\n  RETURN i"))
	if err != nil {
		t.Fatal(err)
	}
	defer plan.Close()
	session, err := plan.NewDebugSession(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	defer session.Close()

	if _, err := session.Start(context.Background()); err != nil {
		t.Fatal(err)
	}
	first, err := session.Step(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	second, err := session.Step(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if first.Location.Line != 2 || second.Location.Line != 2 {
		t.Fatalf("expected repeated loop stops, got %#v then %#v", first, second)
	}
}
