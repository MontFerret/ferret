package debugger

import (
	"context"
	"errors"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/source"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

func TestSessionCloseLifecycleStates(t *testing.T) {
	tests := []struct {
		startEvent     *vm.DebugExecutionEvent
		name           string
		wantAfterCalls int
		start          bool
		wantCanceled   bool
	}{
		{
			name: "before_start",
		},
		{
			name:           "while_stopped",
			start:          true,
			startEvent:     &vm.DebugExecutionEvent{Reason: vm.DebugStopEntry},
			wantAfterCalls: 1,
			wantCanceled:   true,
		},
		{
			name:           "after_completion",
			start:          true,
			startEvent:     &vm.DebugExecutionEvent{Reason: vm.DebugStopCompleted, Result: &vm.Result{}},
			wantAfterCalls: 1,
		},
		{
			name:           "after_runtime_error",
			start:          true,
			startEvent:     &vm.DebugExecutionEvent{Reason: vm.DebugStopRuntimeError, Error: errors.New("runtime failed")},
			wantAfterCalls: 1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			src := source.New("close.fql", "RETURN 1")
			point := bytecode.DebugPoint{ID: 1, PC: 1, Span: source.Span{Start: 0, End: 8}, FunctionID: -1}
			if tc.startEvent != nil {
				tc.startEvent.Point = &point
			}
			execution := &fakeExecution{startEvent: tc.startEvent, status: vm.DebugExecutionNew}
			services := &fakeSessionServices{}
			session, err := NewSession(Config{
				Execution:   execution,
				Values:      &fakeValueAccess{inner: vm.NewDebugValueAccess()},
				Services:    services,
				Source:      src,
				DebugPoints: []bytecode.DebugPoint{point},
			})
			if err != nil {
				t.Fatal(err)
			}

			if tc.start {
				if _, err := session.Start(context.Background()); err != nil {
					t.Fatal(err)
				}
			}
			if err := session.Close(); err != nil {
				t.Fatal(err)
			}
			if err := session.Close(); err != nil {
				t.Fatal(err)
			}
			if services.afterCalls != tc.wantAfterCalls {
				t.Fatalf("unexpected after-run calls: got %d, want %d", services.afterCalls, tc.wantAfterCalls)
			}
			if tc.wantCanceled && !errors.Is(services.afterRunErr, context.Canceled) {
				t.Fatalf("expected cancellation after-run error, got %v", services.afterRunErr)
			}
			if !execution.closed || !services.closed {
				t.Fatal("expected execution and services to close")
			}
			if _, err := session.Step(context.Background()); err == nil || !errors.Is(err, &StateError{}) {
				t.Fatalf("expected command rejection after close, got %v", err)
			}
		})
	}
}
