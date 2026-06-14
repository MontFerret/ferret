package ferret

import (
	"context"
	"strings"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/source"
)

func TestDebugSessionBreakpointsDistinguishUDFBodyAndCallSite(t *testing.T) {
	engine, err := New()
	if err != nil {
		t.Fatal(err)
	}
	defer engine.Close()

	query := `LET seed = 1
FUNC add(a) (
  LET b = a + 1
  RETURN b
)
LET value = add(seed)
RETURN value`
	plan, err := engine.CompileDebug(context.Background(), source.New("udf-breakpoints.fql", query))
	if err != nil {
		t.Fatal(err)
	}
	defer plan.Close()

	session, err := plan.NewDebugSession(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	defer session.Close()

	beforeBody, err := session.SetBreakpoint("udf-breakpoints.fql", 2)
	if err != nil {
		t.Fatal(err)
	}
	body, err := session.SetBreakpoint("udf-breakpoints.fql", 3)
	if err != nil {
		t.Fatal(err)
	}
	callSite, err := session.SetBreakpoint("udf-breakpoints.fql", 6)
	if err != nil {
		t.Fatal(err)
	}

	if !beforeBody.Bound || beforeBody.Line != 3 || beforeBody.PointID != body.PointID || beforeBody.FunctionID != body.FunctionID {
		t.Fatalf("non-executable UDF declaration did not bind to its body: before=%#v body=%#v", beforeBody, body)
	}
	if !body.Bound || body.FunctionID < 0 {
		t.Fatalf("expected UDF body breakpoint identity: %#v", body)
	}
	if !callSite.Bound || callSite.FunctionID != -1 || callSite.PointID == body.PointID {
		t.Fatalf("expected distinct caller breakpoint identity: body=%#v call=%#v", body, callSite)
	}

	event, err := session.Start(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if event.Location.Line != 1 {
		t.Fatalf("unexpected entry: %#v", event)
	}
	event, err = session.Continue(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if event.Reason != DebugReasonBreakpoint || event.Location.Line != 6 {
		t.Fatalf("expected caller breakpoint, got %#v", event)
	}
	event, err = session.Continue(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if event.Reason != DebugReasonBreakpoint || event.Location.Line != 3 || event.Depth != 1 {
		t.Fatalf("expected UDF body breakpoint, got %#v", event)
	}
}

func TestDebugSessionBreakpointsDistinguishMultipleUDFs(t *testing.T) {
	engine, err := New()
	if err != nil {
		t.Fatal(err)
	}
	defer engine.Close()

	query := `FUNC first() (
  RETURN 1
)
FUNC second() (
  RETURN 2
)
LET a = first()
LET b = second()
RETURN a + b`
	plan, err := engine.CompileDebug(context.Background(), source.New("multiple-udfs.fql", query))
	if err != nil {
		t.Fatal(err)
	}
	defer plan.Close()

	session, err := plan.NewDebugSession(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	defer session.Close()

	first, err := session.SetBreakpoint("multiple-udfs.fql", 2)
	if err != nil {
		t.Fatal(err)
	}
	second, err := session.SetBreakpoint("multiple-udfs.fql", 5)
	if err != nil {
		t.Fatal(err)
	}
	if !first.Bound || !second.Bound || first.PointID == second.PointID || first.FunctionID == second.FunctionID {
		t.Fatalf("UDF breakpoints were not function-distinct: first=%#v second=%#v", first, second)
	}

	event, err := session.Start(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if event.Location.Line != 7 {
		t.Fatalf("unexpected entry: %#v", event)
	}
	event, err = session.Continue(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if event.Reason != DebugReasonBreakpoint || event.Location.Line != 2 {
		t.Fatalf("expected first UDF breakpoint, got %#v", event)
	}
	event, err = session.Continue(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if event.Reason != DebugReasonBreakpoint || event.Location.Line != 5 {
		t.Fatalf("expected second UDF breakpoint, got %#v", event)
	}
}

func TestDebugSessionNextInFunctionBindingStaysWithinUDFBoundaries(t *testing.T) {
	engine, err := New()
	if err != nil {
		t.Fatal(err)
	}
	defer engine.Close()

	query := `LET seed = 1
FUNC add(a) (
  LET b = a + 1

  RETURN b
)

RETURN add(seed)`
	plan, err := engine.CompileDebug(context.Background(), source.New("udf-binding.fql", query))
	if err != nil {
		t.Fatal(err)
	}
	defer plan.Close()

	session, err := plan.NewDebugSession(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	defer session.Close()

	inside, err := session.SetBreakpointAt(
		DebugSourceLocation{File: "udf-binding.fql", Line: 4},
		DebugBreakpointOptions{BindingMode: DebugBreakpointBindNextExecutableInFunction},
	)
	if err != nil {
		t.Fatal(err)
	}
	if !inside.Bound || inside.Line != 5 || inside.FunctionID < 0 {
		t.Fatalf("expected blank line inside UDF to bind within the UDF: %#v", inside)
	}

	before, err := session.SetBreakpointAt(
		DebugSourceLocation{File: "udf-binding.fql", Line: 2},
		DebugBreakpointOptions{BindingMode: DebugBreakpointBindNextExecutableInFunction},
	)
	if err != nil {
		t.Fatal(err)
	}
	if before.Bound {
		t.Fatalf("function-scoped binding entered a UDF from its declaration boundary: %#v", before)
	}

	after, err := session.SetBreakpointAt(
		DebugSourceLocation{File: "udf-binding.fql", Line: 7},
		DebugBreakpointOptions{BindingMode: DebugBreakpointBindNextExecutableInFunction},
	)
	if err != nil {
		t.Fatal(err)
	}
	if after.Bound {
		t.Fatalf("function-scoped binding left a UDF at its trailing boundary: %#v", after)
	}
}

func TestDebugSessionUDFFramesLocalsAndRuntimeErrorLocation(t *testing.T) {
	engine, err := New()
	if err != nil {
		t.Fatal(err)
	}
	defer engine.Close()

	query := `LET x = 1
FUNC add(a) (
  LET b = a + 1
  RETURN b / 0
)
FUNC unrelated() (
  RETURN 99
)
LET y = add(x)
RETURN y`
	plan, err := engine.CompileDebug(context.Background(), source.New("udf-error.fql", query))
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
	if event.Location.Line != 1 {
		t.Fatalf("unexpected entry: %#v", event)
	}
	event, err = session.Step(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if event.Location.Line != 9 {
		t.Fatalf("unexpected caller step: %#v", event)
	}
	event, err = session.Step(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if event.Location.Line != 3 || event.Depth != 1 {
		t.Fatalf("expected UDF entry, got %#v", event)
	}

	frames, err := session.Frames()
	if err != nil {
		t.Fatal(err)
	}
	if len(frames) != 2 ||
		frames[0].Name != "add" || frames[0].Location.Line != 3 ||
		frames[1].Name != "<main>" || frames[1].Location.Line != 9 {
		t.Fatalf("unexpected UDF frames: %#v", frames)
	}
	locals, err := session.Locals()
	if err != nil {
		t.Fatal(err)
	}
	if len(locals) != 1 || locals[0].Name != "a" || locals[0].Value.Display != "1" {
		t.Fatalf("unexpected UDF entry locals: %#v", locals)
	}

	event, err = session.Step(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if event.Location.Line != 4 {
		t.Fatalf("unexpected UDF return step: %#v", event)
	}
	locals, err = session.Locals()
	if err != nil {
		t.Fatal(err)
	}
	values := make(map[string]string, len(locals))
	for _, local := range locals {
		values[local.Name] = local.Value.Display
	}
	if len(locals) != 2 || values["a"] != "1" || values["b"] != "2" {
		t.Fatalf("unexpected UDF return locals: %#v", locals)
	}

	event, err = session.Continue(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if event.Reason != DebugReasonRuntimeError || event.Location.Line != 4 || event.Error == nil {
		t.Fatalf("unexpected UDF runtime error: %#v", event)
	}
	formatted := diagnostics.Format(event.Error)
	if !strings.Contains(formatted, "udf-error.fql:4") || !strings.Contains(formatted, "RETURN b / 0") {
		t.Fatalf("expected UDF error location, got:\n%s", formatted)
	}
}
