package debugger

import (
	"context"
	"errors"
	"testing"

	apidebugger "github.com/MontFerret/api/debugger"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/source"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

func TestSessionUsesInterfacesForBreakpointsEvaluationAndLifecycle(t *testing.T) {
	src := source.New("debug.fql", "LET x = 1 RETURN x")
	point := bytecode.DebugPoint{ID: 9, PC: 3, Span: source.Span{Start: 0, End: 3}, FunctionID: bytecode.NoFunction}
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

	breakpoint, err := session.SetBreakpoint(source.Location{SourceName: "debug.fql", Position: source.Position{Line: 1}})
	if err != nil {
		t.Fatal(err)
	}

	if !breakpoint.Bound || int(breakpoint.PointID) != int(point.ID) || int(breakpoint.FunctionID) != int(point.FunctionID) || breakpoint.RequestedLocation.Column != 0 {
		t.Fatalf("expected bound breakpoint: %#v", breakpoint)
	}

	if got := session.Breakpoints(); len(got) != 1 || int(got[0].PointID) != int(point.ID) || int(got[0].FunctionID) != int(point.FunctionID) {
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
	if execution.resumeBreakpoints == nil || !execution.resumeBreakpoints(point.PC) ||
		execution.resumeBreakpoints(point.PC+1) {
		t.Fatalf("expected a predicate matching exactly PC %d", point.PC)
	}
	value, err := session.Evaluate(context.Background(), "x + 1")
	if err != nil {
		t.Fatal(err)
	}
	if value.Display != "2" || values.typeCalls == 0 || values.debugInfoCalls == 0 {
		t.Fatalf("unexpected evaluated value: %#v", value)
	}
	if _, err := session.FrameLocals(1); !errors.Is(err, runtime.ErrInvalidOperation) {
		t.Fatalf("expected legacy execution to reject caller inspection, got %v", err)
	}
	value, err = session.EvaluateFrame(context.Background(), 0, "x + 2")
	if err != nil || value.Display != "3" {
		t.Fatalf("unexpected legacy top-frame evaluation: %#v, %v", value, err)
	}
	if err := session.Close(); err != nil {
		t.Fatal(err)
	}
	if !execution.closed || !services.closed {
		t.Fatal("expected execution and services to close")
	}
}

func TestSessionCloseDoesNotRepeatAbortedStartHooks(t *testing.T) {
	execution := &fakeExecution{status: vm.DebugExecutionNew}
	services := &fakeSessionServices{beforeRun: func(ctx context.Context) (context.Context, error) {
		ctx, cancel := context.WithCancel(ctx)
		cancel()

		return ctx, nil
	}}
	session, err := NewSession(Config{
		Execution: execution,
		Values:    vm.NewDebugValueAccess(),
		Services:  services,
		Source:    source.NewAnonymous("RETURN 1"),
	})
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() { _ = session.Close() })
	if event, err := session.Start(t.Context()); event != nil || !errors.Is(err, context.Canceled) {
		t.Fatalf("aborted start: event=%+v err=%v", event, err)
	}

	if execution.Status() != vm.DebugExecutionNew || services.afterCalls != 1 || !errors.Is(services.afterRunErr, context.Canceled) {
		t.Fatalf("aborted start: status=%v after hooks=%d error=%v", execution.Status(), services.afterCalls, services.afterRunErr)
	}

	for range 2 {
		if err := session.Close(); err != nil {
			t.Fatal(err)
		}
	}

	if services.afterCalls != 1 {
		t.Fatalf("close repeated aborted attempt's hooks: %d", services.afterCalls)
	}
}

func TestSessionCloseReturnsAndCachesExecutionCloseError(t *testing.T) {
	closeErr := errors.New("execution close failed")
	src := source.New("close.fql", "RETURN 1")
	point := bytecode.DebugPoint{PC: 0, Span: source.Span{Start: 0, End: 8}, FunctionID: bytecode.NoFunction}
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
		{ID: 4, PC: 5, Span: source.Span{Start: 0, End: 8}, FunctionID: bytecode.NoFunction},
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

	breakpoint, err := session.SetBreakpoint(source.Location{SourceName: "", Position: source.Position{Line: 1}})
	if err != nil {
		t.Fatal(err)
	}

	if !breakpoint.Bound || breakpoint.PointID != 4 || breakpoint.FunctionID != apidebugger.NoFunction || breakpoint.Location.Line != 1 {
		t.Fatalf("unexpected source-ordered breakpoint: %#v", breakpoint)
	}

	unbound, err := session.SetBreakpoint(source.Location{SourceName: "other.fql", Position: source.Position{Line: 1}})
	if err != nil {
		t.Fatal(err)
	}
	if unbound.Bound {
		t.Fatalf("breakpoint crossed source boundary: %#v", unbound)
	}
}

func TestSessionVariablesExpandNestedCollections(t *testing.T) {
	src := source.New("debug.fql", "RETURN 1")
	point := bytecode.DebugPoint{ID: 1, PC: 0, Span: source.Span{Start: 0, End: 8}, FunctionID: bytecode.NoFunction}
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
	point := bytecode.DebugPoint{ID: 1, PC: 0, Span: source.Span{Start: 0, End: 8}, FunctionID: bytecode.NoFunction}
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
	point := bytecode.DebugPoint{ID: 1, PC: 0, Span: source.Span{Start: 0, End: 8}, FunctionID: bytecode.NoFunction}
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
	point := bytecode.DebugPoint{ID: 1, PC: 0, Span: source.Span{Start: 0, End: 8}, FunctionID: bytecode.NoFunction}
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
	point := bytecode.DebugPoint{ID: 1, PC: 0, Span: source.Span{Start: 0, End: 8}, FunctionID: bytecode.NoFunction}
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

func TestSessionAbortedStartPreservesReferencesUntilRetry(t *testing.T) {
	for _, invalid := range []string{"nil_context", "canceled_context"} {
		t.Run(invalid, func(t *testing.T) {
			execution := &fakeExecution{
				startEvent: &vm.DebugExecutionEvent{Reason: vm.DebugStopEntry},
				locals:     []vm.DebugLocal{{Name: "value", Value: runtime.NewArrayWith(runtime.NewInt(1))}},
				status:     vm.DebugExecutionNew,
			}
			beforeCalls := 0
			services := &fakeSessionServices{beforeRun: func(ctx context.Context) (context.Context, error) {
				beforeCalls++
				if beforeCalls != 1 {
					return ctx, nil
				}

				if invalid == "nil_context" {
					return nil, nil
				}

				canceled, cancel := context.WithCancel(ctx)
				cancel()

				return canceled, nil
			}}
			session, err := NewSession(Config{
				Execution: execution,
				Values:    vm.NewDebugValueAccess(),
				Services:  services,
				Source:    source.NewAnonymous("RETURN 1"),
			})
			if err != nil {
				t.Fatal(err)
			}
			t.Cleanup(func() { _ = session.Close() })

			locals, err := session.Locals()
			if err != nil {
				t.Fatal(err)
			}
			reference := locals[0].Value.Reference
			expectedErr := runtime.ErrInvalidArgument
			if invalid == "canceled_context" {
				expectedErr = context.Canceled
			}

			if event, err := session.Start(t.Context()); event != nil || !errors.Is(err, expectedErr) {
				t.Fatalf("aborted start: event=%+v err=%v", event, err)
			}

			if execution.Status() != vm.DebugExecutionNew || services.afterCalls != 1 || !errors.Is(services.afterRunErr, expectedErr) {
				t.Fatalf("aborted attempt: status=%v hooks=%d err=%v", execution.Status(), services.afterCalls, services.afterRunErr)
			}

			if _, err := session.Variables(reference); err != nil {
				t.Fatalf("aborted start invalidated reference: %v", err)
			}

			if event, err := session.Start(t.Context()); err != nil || event.Reason != ReasonEntry {
				t.Fatalf("retry: event=%+v err=%v", event, err)
			}

			if _, err := session.Variables(reference); !errors.Is(err, runtime.ErrNotFound) {
				t.Fatalf("successful start retained stale reference: %v", err)
			}

			if services.afterCalls != 1 {
				t.Fatalf("retry settled hooks too early: %d", services.afterCalls)
			}

			for range 2 {
				if err := session.Close(); err != nil {
					t.Fatal(err)
				}
			}

			if beforeCalls != 2 || services.afterCalls != 2 || !errors.Is(services.afterRunErr, context.Canceled) {
				t.Fatalf("retry hook pairing: before=%d after=%d err=%v", beforeCalls, services.afterCalls, services.afterRunErr)
			}
		})
	}
}

func BenchmarkSessionContinueThroughExecutionInterface(b *testing.B) {
	src := source.New("debug.fql", "RETURN 1")
	point := bytecode.DebugPoint{ID: 7, PC: 0, Span: source.Span{Start: 0, End: 6}, FunctionID: bytecode.NoFunction}
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
