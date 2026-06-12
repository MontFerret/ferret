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
	point := bytecode.DebugPoint{PC: 3, Span: source.Span{Start: 0, End: 3}, FunctionID: -1}
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
	if !breakpoint.Bound {
		t.Fatalf("expected bound breakpoint: %#v", breakpoint)
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

func BenchmarkSessionContinueThroughExecutionInterface(b *testing.B) {
	src := source.New("debug.fql", "RETURN 1")
	point := bytecode.DebugPoint{PC: 0, Span: source.Span{Start: 0, End: 6}, FunctionID: -1}
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
