package debugger

import (
	"context"
	"errors"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/source"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

func TestSessionUsesInterfacesForBreakpointsEvaluationAndLifecycle(t *testing.T) {
	src := source.New("debug.fql", "LET x = 1 RETURN x")
	point := bytecode.DebugPoint{ID: 9, PC: 3, Span: source.Span{Start: 0, End: 3}, FunctionID: -1}
	execution := &fakeExecution{
		startEvent:  &vm.DebugExecutionEvent{Reason: vm.DebugStopEntry, Point: &point},
		resumeEvent: &vm.DebugExecutionEvent{Reason: vm.DebugStopStep, Point: &point},
		locals:      []vm.DebugLocal{{Name: "x", Value: runtime.NewInt(1)}},
		params:      runtime.NewParams(),
		status:      vm.DebugExecutionNew,
	}
	services := &fakeSessionServices{}
	values := &fakeValueAccess{inner: vm.NewDebugValueAccess()}
	session, err := NewSession(Config{
		Execution:   execution,
		Values:      values,
		Services:    services,
		Source:      src,
		DebugPoints: []bytecode.DebugPoint{point},
	})
	if err != nil {
		t.Fatal(err)
	}

	breakpoint, err := session.SetBreakpoint("debug.fql", 1)
	if err != nil {
		t.Fatal(err)
	}
	if !breakpoint.Bound || breakpoint.PointID != point.ID || breakpoint.FunctionID != point.FunctionID || breakpoint.RequestedColumn != 0 {
		t.Fatalf("expected bound breakpoint: %#v", breakpoint)
	}
	if got := session.Breakpoints(); len(got) != 1 || got[0].PointID != point.ID || got[0].FunctionID != point.FunctionID {
		t.Fatalf("breakpoint snapshot lost bound identity: %#v", got)
	}
	if _, err := session.Start(context.Background()); err != nil {
		t.Fatal(err)
	}
	if _, err := session.Start(context.Background()); err == nil || !errors.Is(err, &StateError{}) {
		t.Fatalf("expected typed state error, got %v", err)
	}
	if _, err := session.Continue(context.Background()); err != nil {
		t.Fatal(err)
	}
	if len(execution.resumeBreakpoints) != 1 {
		t.Fatalf("expected one exact breakpoint PC, got %#v", execution.resumeBreakpoints)
	}
	if _, exists := execution.resumeBreakpoints[point.PC]; !exists {
		t.Fatalf("expected point id %d to resolve to pc %d, got %#v", point.ID, point.PC, execution.resumeBreakpoints)
	}
	value, err := session.Evaluate(context.Background(), "x + 1")
	if err != nil {
		t.Fatal(err)
	}
	if value.Display != "2" || values.typeCalls == 0 || values.debugInfoCalls == 0 {
		t.Fatalf("unexpected evaluated value: %#v", value)
	}
	if err := session.Close(); err != nil {
		t.Fatal(err)
	}
	if !execution.closed || !services.closed {
		t.Fatal("expected execution and services to close")
	}
}

func TestSessionCloseReturnsAndCachesExecutionCloseError(t *testing.T) {
	closeErr := errors.New("execution close failed")
	src := source.New("close.fql", "RETURN 1")
	point := bytecode.DebugPoint{PC: 0, Span: source.Span{Start: 0, End: 8}, FunctionID: -1}
	execution := &fakeExecution{
		closeErr: closeErr,
		status:   vm.DebugExecutionNew,
	}
	session, err := NewSession(Config{
		Execution:   execution,
		Values:      &fakeValueAccess{inner: vm.NewDebugValueAccess()},
		Services:    &fakeSessionServices{},
		Source:      src,
		DebugPoints: []bytecode.DebugPoint{point},
	})
	if err != nil {
		t.Fatal(err)
	}

	if err := session.Close(); !errors.Is(err, closeErr) {
		t.Fatalf("expected execution close error, got %v", err)
	}
	if err := session.Close(); !errors.Is(err, closeErr) {
		t.Fatalf("expected cached execution close error, got %v", err)
	}
}

func TestSessionBreakpointBindingUsesSourceOrderAndStableTieBreaks(t *testing.T) {
	src := source.New("binding.fql", "RETURN 1\nRETURN 2")
	points := []bytecode.DebugPoint{
		{ID: 8, PC: 2, Span: source.Span{Start: 9, End: 17}, FunctionID: 1},
		{ID: 4, PC: 5, Span: source.Span{Start: 0, End: 8}, FunctionID: -1},
		{ID: 3, PC: 7, Span: source.Span{Start: 0, End: 8}, FunctionID: 0},
	}
	session, err := NewSession(Config{
		Execution:   &fakeExecution{status: vm.DebugExecutionNew},
		Values:      &fakeValueAccess{inner: vm.NewDebugValueAccess()},
		Services:    &fakeSessionServices{},
		Source:      src,
		DebugPoints: points,
	})
	if err != nil {
		t.Fatal(err)
	}
	defer session.Close()

	breakpoint, err := session.SetBreakpoint("", 1)
	if err != nil {
		t.Fatal(err)
	}
	if !breakpoint.Bound || breakpoint.PointID != 4 || breakpoint.FunctionID != -1 || breakpoint.Line != 1 {
		t.Fatalf("unexpected source-ordered breakpoint: %#v", breakpoint)
	}

	unbound, err := session.SetBreakpoint("other.fql", 1)
	if err != nil {
		t.Fatal(err)
	}
	if unbound.Bound {
		t.Fatalf("breakpoint crossed source boundary: %#v", unbound)
	}
}

func TestSessionVariablesExpandNestedCollections(t *testing.T) {
	src := source.New("debug.fql", "RETURN 1")
	point := bytecode.DebugPoint{ID: 1, PC: 0, Span: source.Span{Start: 0, End: 8}, FunctionID: -1}
	execution := &fakeExecution{
		startEvent: &vm.DebugExecutionEvent{Reason: vm.DebugStopEntry, Point: &point},
		locals: []vm.DebugLocal{{
			Name: "value",
			Value: runtime.NewObjectWith(map[string]runtime.Value{
				"nested": runtime.NewArrayWith(
					runtime.NewInt(1),
					runtime.NewObjectWith(map[string]runtime.Value{
						"answer": runtime.NewInt(42),
					}),
				),
			}),
		}},
		status: vm.DebugExecutionNew,
	}
	session, err := NewSession(Config{
		Execution:   execution,
		Values:      &fakeValueAccess{inner: vm.NewDebugValueAccess()},
		Services:    &fakeSessionServices{},
		Source:      src,
		DebugPoints: []bytecode.DebugPoint{point},
	})
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
	if len(locals) != 1 || !locals[0].Value.Reference.Valid() {
		t.Fatalf("expected expandable local, got %#v", locals)
	}

	children, err := session.Variables(locals[0].Value.Reference)
	if err != nil {
		t.Fatal(err)
	}
	if len(children) != 1 || children[0].Name != "nested" || !children[0].Value.Reference.Valid() {
		t.Fatalf("unexpected object children: %#v", children)
	}

	items, err := session.Variables(children[0].Value.Reference)
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 2 || items[0].Name != "0" || items[0].Value.Display != "1" || items[1].Name != "1" || !items[1].Value.Reference.Valid() {
		t.Fatalf("unexpected array children: %#v", items)
	}

	nested, err := session.Variables(items[1].Value.Reference)
	if err != nil {
		t.Fatal(err)
	}
	if len(nested) != 1 || nested[0].Name != "answer" || nested[0].Value.Display != "42" {
		t.Fatalf("unexpected nested object children: %#v", nested)
	}
}

func TestSessionVariablesStaySummarizedWhenCollectionExceedsBounds(t *testing.T) {
	src := source.New("debug.fql", "RETURN 1")
	point := bytecode.DebugPoint{ID: 1, PC: 0, Span: source.Span{Start: 0, End: 8}, FunctionID: -1}
	execution := &fakeExecution{
		startEvent: &vm.DebugExecutionEvent{Reason: vm.DebugStopEntry, Point: &point},
		locals: []vm.DebugLocal{{
			Name:  "value",
			Value: runtime.NewArrayWith(runtime.NewInt(1), runtime.NewInt(2)),
		}},
		status: vm.DebugExecutionNew,
	}
	session, err := NewSession(Config{
		Execution:   execution,
		Values:      &fakeValueAccess{inner: vm.NewDebugValueAccess()},
		Services:    &fakeSessionServices{},
		Source:      src,
		DebugPoints: []bytecode.DebugPoint{point},
		Format:      FormatOptions{MaxDepth: 3, MaxItems: 1, MaxBytes: 1024},
	})
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
	if len(locals) != 1 || locals[0].Value.Reference.Valid() || locals[0].Value.Display != "Array(2)" {
		t.Fatalf("expected bounded summary without reference, got %#v", locals)
	}

	if _, err := session.Variables(0); !errors.Is(err, runtime.ErrInvalidArgument) {
		t.Fatalf("expected invalid reference error, got %v", err)
	}
}

func TestSessionVariablesInvalidateReferencesAfterResume(t *testing.T) {
	src := source.New("debug.fql", "RETURN 1")
	point := bytecode.DebugPoint{ID: 1, PC: 0, Span: source.Span{Start: 0, End: 8}, FunctionID: -1}
	execution := &fakeExecution{
		startEvent:  &vm.DebugExecutionEvent{Reason: vm.DebugStopEntry, Point: &point},
		resumeEvent: &vm.DebugExecutionEvent{Reason: vm.DebugStopStep, Point: &point},
		locals: []vm.DebugLocal{{
			Name:  "value",
			Value: runtime.NewArrayWith(runtime.NewInt(1)),
		}},
		status: vm.DebugExecutionNew,
	}
	session, err := NewSession(Config{
		Execution:   execution,
		Values:      &fakeValueAccess{inner: vm.NewDebugValueAccess()},
		Services:    &fakeSessionServices{},
		Source:      src,
		DebugPoints: []bytecode.DebugPoint{point},
	})
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
	reference := locals[0].Value.Reference
	if !reference.Valid() {
		t.Fatalf("expected expandable reference, got %#v", locals[0].Value)
	}

	if _, err := session.Continue(context.Background()); err != nil {
		t.Fatal(err)
	}

	if _, err := session.Variables(reference); !errors.Is(err, runtime.ErrNotFound) {
		t.Fatalf("expected stale reference to be rejected, got %v", err)
	}
}

func TestSessionVariablesInvalidateReferencesAtStart(t *testing.T) {
	src := source.New("debug.fql", "RETURN 1")
	point := bytecode.DebugPoint{ID: 1, PC: 0, Span: source.Span{Start: 0, End: 8}, FunctionID: -1}
	execution := &fakeExecution{
		startEvent: &vm.DebugExecutionEvent{Reason: vm.DebugStopEntry, Point: &point},
		locals: []vm.DebugLocal{{
			Name:  "value",
			Value: runtime.NewArrayWith(runtime.NewInt(1)),
		}},
		status: vm.DebugExecutionNew,
	}
	session, err := NewSession(Config{
		Execution:   execution,
		Values:      &fakeValueAccess{inner: vm.NewDebugValueAccess()},
		Services:    &fakeSessionServices{},
		Source:      src,
		DebugPoints: []bytecode.DebugPoint{point},
	})
	if err != nil {
		t.Fatal(err)
	}
	defer session.Close()

	locals, err := session.Locals()
	if err != nil {
		t.Fatal(err)
	}
	reference := locals[0].Value.Reference
	if !reference.Valid() {
		t.Fatalf("expected expandable reference, got %#v", locals[0].Value)
	}

	if _, err := session.Start(context.Background()); err != nil {
		t.Fatal(err)
	}
	if _, err := session.Variables(reference); !errors.Is(err, runtime.ErrNotFound) {
		t.Fatalf("expected pre-start reference to be rejected, got %v", err)
	}
}

func TestSessionVariablesRejectUnknownCollectionKind(t *testing.T) {
	src := source.New("debug.fql", "RETURN 1")
	point := bytecode.DebugPoint{ID: 1, PC: 0, Span: source.Span{Start: 0, End: 8}, FunctionID: -1}
	value := runtime.NewArrayWith(runtime.NewInt(1))
	execution := &fakeExecution{
		startEvent: &vm.DebugExecutionEvent{Reason: vm.DebugStopEntry, Point: &point},
		locals:     []vm.DebugLocal{{Name: "value", Value: value}},
		status:     vm.DebugExecutionNew,
	}
	values := &fakeValueAccess{
		inner: vm.NewDebugValueAccess(),
		inspect: func(runtime.Value, int) (vm.DebugValueInspection, bool) {
			return vm.DebugValueInspection{
				Kind:     vm.DebugValueKind(255),
				Length:   1,
				Items:    []vm.DebugValueItem{{Value: runtime.NewInt(1)}},
				Complete: true,
			}, true
		},
	}
	session, err := NewSession(Config{
		Execution:   execution,
		Values:      values,
		Services:    &fakeSessionServices{},
		Source:      src,
		DebugPoints: []bytecode.DebugPoint{point},
	})
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
	if len(locals) != 1 || locals[0].Value.Reference.Valid() {
		t.Fatalf("expected unknown collection kind to remain non-expandable, got %#v", locals)
	}
}

func BenchmarkSessionContinueThroughExecutionInterface(b *testing.B) {
	src := source.New("debug.fql", "RETURN 1")
	point := bytecode.DebugPoint{ID: 7, PC: 0, Span: source.Span{Start: 0, End: 6}, FunctionID: -1}
	execution := &fakeExecution{
		startEvent:  &vm.DebugExecutionEvent{Reason: vm.DebugStopEntry, Point: &point},
		resumeEvent: &vm.DebugExecutionEvent{Reason: vm.DebugStopStep, Point: &point},
		params:      runtime.NewParams(),
		status:      vm.DebugExecutionNew,
	}
	session, err := NewSession(Config{
		Execution:   execution,
		Values:      &fakeValueAccess{inner: vm.NewDebugValueAccess()},
		Services:    &fakeSessionServices{},
		Source:      src,
		DebugPoints: []bytecode.DebugPoint{point},
	})
	if err != nil {
		b.Fatal(err)
	}
	if _, err := session.Start(context.Background()); err != nil {
		b.Fatal(err)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := session.Continue(context.Background()); err != nil {
			b.Fatal(err)
		}
	}
}
