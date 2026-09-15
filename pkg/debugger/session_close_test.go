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
			point := bytecode.DebugPoint{ID: 1, PC: 1, Span: source.Span{Start: 0, End: 8}, FunctionID: bytecode.NoFunction}
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
			if _, err := session.StepIn(context.Background()); err == nil || !errors.Is(err, &StateError{}) {
				t.Fatalf("expected command rejection after close, got %v", err)
			}
		})
	}
}

func TestSessionNilAndClosedReceivers(t *testing.T) {
	for _, tc := range []struct {
		session *Session
		name    string
	}{
		{name: "nil"},
		{name: "closed_zero_value", session: &Session{}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			session := tc.session
			for range 2 {
				if err := session.Close(); err != nil {
					t.Fatal(err)
				}
			}

			listed := session.Breakpoints()
			if len(listed) != 0 || (listed == nil) != (session == nil) {
				t.Fatalf("unexpected empty listing: %#v", listed)
			}

			commands := []struct {
				call func(context.Context) (*Event, error)
				name string
			}{
				{name: "start", call: session.Start},
				{name: "continue", call: session.Continue},
				{name: "step_in", call: session.StepIn},
				{name: "step_over", call: session.StepOver},
				{name: "step_out", call: session.StepOut},
			}
			for _, command := range commands {
				if _, err := command.call(t.Context()); !errors.Is(err, &StateError{}) {
					t.Fatalf("%s: expected closed error, got %v", command.name, err)
				}

				if _, err := command.call(nil); !errors.Is(err, runtime.ErrInvalidArgument) {
					t.Fatalf("%s: context validation lost priority: %v", command.name, err)
				}
			}

			if err := session.Pause(); !errors.Is(err, &StateError{}) {
				t.Fatalf("pause: %v", err)
			}

			if _, err := session.Frames(); !errors.Is(err, &StateError{}) {
				t.Fatalf("frames: %v", err)
			}

			if _, err := session.FrameLocals(0); !errors.Is(err, &StateError{}) {
				t.Fatalf("locals: %v", err)
			}

			if _, err := session.Variables(1); !errors.Is(err, &StateError{}) {
				t.Fatalf("variables: %v", err)
			}

			if _, err := session.Evaluate(t.Context(), "1"); !errors.Is(err, &StateError{}) {
				t.Fatalf("evaluate: %v", err)
			}

			if _, err := session.ReplaceBreakpoints(t.Context(), "", nil); !errors.Is(err, &StateError{}) {
				t.Fatalf("replace breakpoints: %v", err)
			}

			if _, err := session.SetBreakpoint(source.Location{}); !errors.Is(err, &StateError{}) {
				t.Fatalf("add breakpoint: %v", err)
			}

			if err := session.DeleteBreakpoint(1); !errors.Is(err, &StateError{}) {
				t.Fatalf("delete breakpoint: %v", err)
			}
		})
	}
}
