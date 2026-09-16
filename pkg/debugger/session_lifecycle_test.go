package debugger

import (
	"context"
	"errors"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/source"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

func TestSessionDirectTerminalErrorsSettleAfterRun(t *testing.T) {
	runErr := errors.New("execution failed")
	hookErr := errors.New("hook failed")
	executionCloseErr := errors.New("execution cleanup failed")
	servicesCloseErr := errors.New("services cleanup failed")

	for _, command := range []string{"start", "resume"} {
		for _, tc := range []struct {
			hookErr           error
			executionCloseErr error
			servicesCloseErr  error
			name              string
			terminalState     string
			status            vm.DebugExecutionStatus
		}{
			{name: "terminated", status: vm.DebugExecutionTerminated, terminalState: "terminated"},
			{name: "completed", status: vm.DebugExecutionCompleted, terminalState: "completed"},
			{name: "closed", status: vm.DebugExecutionClosed, terminalState: "terminated"},
			{name: "hook_failure", status: vm.DebugExecutionTerminated, terminalState: "terminated", hookErr: hookErr},
			{name: "cleanup_failures", status: vm.DebugExecutionTerminated, terminalState: "terminated", executionCloseErr: executionCloseErr, servicesCloseErr: servicesCloseErr},
		} {
			t.Run(command+"/"+tc.name, func(t *testing.T) {
				execution := &fakeExecution{
					startEvent:  &vm.DebugExecutionEvent{Reason: vm.DebugStopEntry},
					errorStatus: tc.status,
					closeErr:    tc.executionCloseErr,
				}
				services := &fakeSessionServices{closeErr: tc.servicesCloseErr}
				session := newLifecycleTestSession(t, execution, services)
				services.afterRun = func(_ context.Context, err error) error {
					if err != runErr {
						t.Errorf("AfterRun error = %v, want original execution error", err)
					}

					select {
					case <-session.lifecycle.terminalDone:
					default:
						t.Error("AfterRun ran before terminal commitment")
					}

					if session.lifecycle.terminalState != tc.terminalState {
						t.Errorf("terminal state = %q, want %q", session.lifecycle.terminalState, tc.terminalState)
					}

					return tc.hookErr
				}

				call := session.Start
				if command == "resume" {
					if _, err := session.Start(t.Context()); err != nil {
						t.Fatal(err)
					}

					execution.resumeErr = runErr
					call = session.Continue
				} else {
					execution.startErr = runErr
				}

				event, err := call(t.Context())
				if event != nil || !errors.Is(err, runErr) {
					t.Fatalf("direct failure: event=%+v err=%v", event, err)
				}

				if tc.hookErr == nil && err != runErr {
					t.Fatalf("execution error identity changed: %v", err)
				}

				if tc.hookErr != nil && !errors.Is(err, tc.hookErr) {
					t.Fatalf("lost hook failure: %v", err)
				}

				if services.afterCalls != 1 || services.afterRunErr != runErr {
					t.Fatalf("AfterRun not settled immediately: calls=%d err=%v", services.afterCalls, services.afterRunErr)
				}

				if execution.closeCalls != 0 || services.closeCalls != 0 {
					t.Fatal("terminal failure released resources before Close")
				}

				var firstCloseErr error
				for i := range 2 {
					err := session.Close()
					if tc.executionCloseErr == nil && err != nil {
						t.Fatalf("clean Close returned an error: %v", err)
					}

					if tc.executionCloseErr != nil && (!errors.Is(err, tc.executionCloseErr) || !errors.Is(err, tc.servicesCloseErr)) {
						t.Fatalf("Close lost cleanup failures: %v", err)
					}

					if errors.Is(err, runErr) || errors.Is(err, hookErr) || errors.Is(err, context.Canceled) {
						t.Fatalf("Close retained a run failure: %v", err)
					}

					if i == 0 {
						firstCloseErr = err
					} else if err != firstCloseErr {
						t.Fatal("Close did not return its cached cleanup error")
					}
				}

				if services.afterCalls != 1 || services.afterRunErr != runErr || execution.closeCalls != 1 || services.closeCalls != 1 {
					t.Fatalf("Close repeated lifecycle work: after=%d err=%v execution=%d services=%d", services.afterCalls, services.afterRunErr, execution.closeCalls, services.closeCalls)
				}
			})
		}
	}
}

func TestSessionDirectNonterminalErrorsDoNotSettleAfterRun(t *testing.T) {
	for _, command := range []string{"start", "resume"} {
		t.Run(command, func(t *testing.T) {
			execution := &fakeExecution{startEvent: &vm.DebugExecutionEvent{Reason: vm.DebugStopEntry}}
			services := &fakeSessionServices{}
			session := newLifecycleTestSession(t, execution, services)
			call := session.Start
			if command == "resume" {
				if _, err := session.Start(t.Context()); err != nil {
					t.Fatal(err)
				}

				execution.resumeErr = context.Canceled
				execution.errorStatus = vm.DebugExecutionPaused
				call = session.Continue
			} else {
				execution.startErr = context.Canceled
				execution.errorStatus = vm.DebugExecutionNew
			}

			if event, err := call(t.Context()); event != nil || err != context.Canceled {
				t.Fatalf("direct admission failure: event=%+v err=%v", event, err)
			}

			select {
			case <-session.lifecycle.terminalDone:
				t.Fatal("nonterminal failure committed terminal state")
			default:
			}

			if services.afterCalls != 0 || execution.closeCalls != 0 || services.closeCalls != 0 {
				t.Fatal("nonterminal failure settled or closed the run")
			}

			for range 2 {
				if err := session.Close(); err != nil {
					t.Fatal(err)
				}
			}

			if services.afterCalls != 1 || services.afterRunErr != context.Canceled || execution.closeCalls != 1 || services.closeCalls != 1 {
				t.Fatalf("Close did not settle cancellation once: after=%d err=%v execution=%d services=%d", services.afterCalls, services.afterRunErr, execution.closeCalls, services.closeCalls)
			}
		})
	}
}

func TestSessionRuntimeErrorStopSettlesHooksWithoutTermination(t *testing.T) {
	for _, finish := range []string{"resume", "close"} {
		t.Run(finish, func(t *testing.T) {
			runErr := errors.New("runtime failed")
			execution := &fakeExecution{
				startEvent:  &vm.DebugExecutionEvent{Reason: vm.DebugStopEntry},
				resumeEvent: &vm.DebugExecutionEvent{Reason: vm.DebugStopRuntimeError, Error: runErr},
				locals:      []vm.DebugLocal{{Name: "value", Value: runtime.NewInt(7)}},
				frames:      []vm.DebugFrame{{Name: "main"}},
			}
			services := &fakeSessionServices{}
			session := newLifecycleTestSession(t, execution, services)
			if _, err := session.Start(t.Context()); err != nil {
				t.Fatal(err)
			}

			event, err := session.Continue(t.Context())
			if err != nil || event == nil || event.Reason != ReasonRuntimeError || event.Error != runErr {
				t.Fatalf("runtime error stop: event=%+v err=%v", event, err)
			}

			select {
			case <-session.lifecycle.terminalDone:
				t.Fatal("runtime error stop committed terminal state")
			default:
			}

			if execution.Status() != vm.DebugExecutionPaused || services.afterCalls != 1 || services.afterRunErr != runErr || execution.closeCalls != 0 || services.closeCalls != 0 {
				t.Fatal("runtime error stop did not retain paused state and settle its hook")
			}

			if frames, err := session.Frames(t.Context()); err != nil || len(frames) != 1 {
				t.Fatalf("frames at runtime error stop: %+v, %v", frames, err)
			}

			if locals, err := session.FrameLocals(t.Context(), 0); err != nil || len(locals) != 1 || locals[0].Value.Display != "7" {
				t.Fatalf("locals at runtime error stop: %+v, %v", locals, err)
			}

			if value, err := session.EvaluateFrame(t.Context(), 0, "value + 1"); err != nil || value.Display != "8" {
				t.Fatalf("evaluation at runtime error stop: %+v, %v", value, err)
			}

			if finish == "resume" {
				execution.resumeEvent = &vm.DebugExecutionEvent{Reason: vm.DebugStopTerminated, Error: runErr}
				event, err := session.Continue(t.Context())
				if err != nil || event == nil || event.Reason != ReasonTerminated || event.Error != runErr {
					t.Fatalf("termination after runtime error stop: %+v, %v", event, err)
				}

				select {
				case <-session.lifecycle.terminalDone:
				default:
					t.Fatal("resume after runtime error stop did not commit termination")
				}
			}

			for range 2 {
				if err := session.Close(); err != nil {
					t.Fatal(err)
				}
			}

			if services.afterCalls != 1 || services.afterRunErr != runErr || execution.closeCalls != 1 || services.closeCalls != 1 {
				t.Fatal("termination or Close repeated runtime-error lifecycle work")
			}
		})
	}
}

func newLifecycleTestSession(t *testing.T, execution vm.DebugExecution, services SessionServices) *Session {
	t.Helper()

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

	return session
}
