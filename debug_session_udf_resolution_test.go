package ferret

import (
	"context"
	"strings"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/source"
)

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
