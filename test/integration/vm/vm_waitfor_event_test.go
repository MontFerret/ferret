package vm_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
	"github.com/MontFerret/ferret/v2/test/spec"
	. "github.com/MontFerret/ferret/v2/test/spec/exec"
	. "github.com/MontFerret/ferret/v2/test/spec/mock"
)

func TestWaitforEvent(t *testing.T) {
	matchFirst := NewObservable([]runtime.Value{
		NewTestEventType("match"),
		NewTestEventType("other"),
	})
	matchSecond := NewObservable([]runtime.Value{
		NewTestEventType("other"),
		NewTestEventType("match"),
	})
	blocking := NewBlockingObservable()

	RunSpecs(t, []spec.Spec{
		Error(`LET obj = {}

WAITFOR EVENT "test" IN obj

RETURN NONE`, "Should compile but return an error during execution because the object does not implement the interface"),
		S(`LET obj = {}

WAITFOR EVENT "test" IN obj ON ERROR RETURN NONE

RETURN 1`, 1, "Statement suppression should continue after WAITFOR EVENT runtime failure"),
		S(`LET obj = {}

LET status = WAITFOR EVENT "test" IN obj TIMEOUT 1ms ON TIMEOUT RETURN "timeout" ON ERROR RETURN "error"

RETURN status`, "error", "WAITFOR EVENT should choose ON ERROR for runtime failures even when ON TIMEOUT is present"),
		S(`LET obj = {}

LET status = (WAITFOR EVENT "test" IN obj TIMEOUT 1ms) ON TIMEOUT RETURN "timeout" ON ERROR RETURN "error"

RETURN status`, "error", "Grouped WAITFOR EVENT should choose ON ERROR for runtime failures"),
		S(`LET obs = @obs

LET evt = WAITFOR EVENT "test" IN obs

RETURN evt.type`, "match", "WAITFOR EVENT should return the received event value").Env(vm.WithParams(map[string]runtime.Value{
			"obs": matchFirst,
		})),
		Fn(`LET obs = @obs
WAITFOR EVENT "test" IN obs WHEN .type == "match"
RETURN 1`, ObservableReturnOneAndReads(matchFirst, 1)).Env(vm.WithParams(map[string]runtime.Value{
			"obs": matchFirst,
		})),
		Fn(`LET obs = @obs
WAITFOR EVENT "test" IN obs WHEN .type == "match"
RETURN 1`, ObservableReturnOneAndReads(matchSecond, 2)).Env(vm.WithParams(map[string]runtime.Value{
			"obs": matchSecond,
		})),
		Fn(`LET obs = @obs
WAITFOR EVENT "test" IN obs WHEN .type != "" WHEN .type == "match"
RETURN 1`, ObservableReturnOneAndReads(matchSecond, 2)).Env(vm.WithParams(map[string]runtime.Value{
			"obs": matchSecond,
		})),
		S(`LET obs = @obs

LET evt = WAITFOR EVENT "test" IN obs WHEN .type == "match"

RETURN evt.type`, "match", "WAITFOR EVENT filter should return the matched event value").Env(vm.WithParams(map[string]runtime.Value{
			"obs": matchSecond,
		})),
		S(`LET obs = @obs

LET evt = WAITFOR EVENT "test" IN obs WHEN .type != "" WHEN .type == "match"

RETURN evt.type`, "match", "WAITFOR EVENT repeated filters should return the matched event value").Env(vm.WithParams(map[string]runtime.Value{
			"obs": matchSecond,
		})),
		Array(`LET obs = @obs
VAR current = 0

FOR WHILE current < 2
	current += 1
	WAITFOR EVENT "test" IN obs WHEN .type == "match"
	RETURN current`, []any{1, 2}, "WAITFOR EVENT should execute as a FOR loop body statement").Env(vm.WithParams(map[string]runtime.Value{
			"obs": matchFirst,
		})),
		S(`LET obs = @obs

LET evt = WAITFOR EVENT "test" IN obs TIMEOUT 1ms ON TIMEOUT RETURN NONE

RETURN evt.type`, "match", "WAITFOR EVENT timeout-aware success should return the event value").Env(vm.WithParams(map[string]runtime.Value{
			"obs": matchFirst,
		})),
		S(`LET obs = @obs

LET evt = WAITFOR EVENT "test" IN obs ON ERROR RETURN NONE

RETURN evt.type`, "match", "WAITFOR EVENT protected recovery success should return the event value").Env(vm.WithParams(map[string]runtime.Value{
			"obs": matchFirst,
		})),
		S(`LET obs = @obs

LET status = WAITFOR EVENT "test" IN obs TIMEOUT 1ms ON TIMEOUT RETURN "timeout" ON ERROR RETURN "error"

RETURN status`, "timeout", "WAITFOR EVENT should choose ON TIMEOUT when the stream times out").Env(vm.WithParams(map[string]runtime.Value{
			"obs": blocking,
		})),
		S(`LET obs = @obs

LET status = (WAITFOR EVENT "test" IN obs TIMEOUT 1ms) ON TIMEOUT RETURN "timeout" ON ERROR RETURN "error"

RETURN status`, "timeout", "Grouped WAITFOR EVENT should choose ON TIMEOUT when the stream times out").Env(vm.WithParams(map[string]runtime.Value{
			"obs": blocking,
		})),
		spec.NewSpec(`LET obs = @obs

RETURN WAITFOR EVENT "test" IN obs TIMEOUT 1ms ON ERROR RETURN "error"`, "WAITFOR EVENT timeout should not be caught by ON ERROR").Expect().ExecError(
			ShouldBeRuntimeError,
			&ExpectedRuntimeError{Contains: []string{runtime.ErrTimeout.Error()}},
		).Env(vm.WithParams(map[string]runtime.Value{
			"obs": blocking,
		})),
		spec.NewSpec(`LET obs = @obs

RETURN (WAITFOR EVENT "test" IN obs TIMEOUT 1ms) ON ERROR RETURN "error"`, "Grouped WAITFOR EVENT timeout should not be caught by ON ERROR").Expect().ExecError(
			ShouldBeRuntimeError,
			&ExpectedRuntimeError{Contains: []string{runtime.ErrTimeout.Error()}},
		).Env(vm.WithParams(map[string]runtime.Value{
			"obs": blocking,
		})),
	})
}

func TestWaitforEventTrigger(t *testing.T) {
	RunSpecFactory(t, func() []spec.Spec {
		armed := NewTriggerObservable()
		timeout := NewTriggerObservable()
		triggerFailure := NewTriggerObservable()
		triggerFailure.FailNextDispatches(1, errors.New("trigger failed"))
		waitFailure := NewTriggerObservable()
		waitFailure.FailReadsWith(errors.New("stream failed"))
		retry := NewTriggerObservable()
		retry.FailNextDispatches(1, errors.New("trigger failed once"))

		return []spec.Spec{
			Fn(`LET target = @target
LET evt = WAITFOR EVENT "test" IN target
	TRIGGER (
		target <- "test"
	)
	TIMEOUT 20ms
RETURN evt.type`, expectTriggerObservable(armed, "test", 1, 1, 1), "WAITFOR EVENT trigger should run after subscription is armed").Env(vm.WithParams(map[string]runtime.Value{
				"target": armed,
			})),
			Fn(`LET target = @target
RETURN WAITFOR EVENT "test" IN target
	TRIGGER ()
	TIMEOUT 1ms
	ON TIMEOUT RETURN "timeout"`, expectTriggerObservable(timeout, "timeout", 1, 0, 1), "WAITFOR EVENT trigger no-op should preserve timeout cleanup").Env(vm.WithParams(map[string]runtime.Value{
				"target": timeout,
			})),
			Fn(`LET target = @target
RETURN WAITFOR EVENT "test" IN target
	TRIGGER (
		target <- "test"
	)
	TIMEOUT 20ms
	ON ERROR RETURN "error"`, expectTriggerObservable(triggerFailure, "error", 1, 1, 1), "WAITFOR EVENT trigger failure should clean up and use ON ERROR").Env(vm.WithParams(map[string]runtime.Value{
				"target": triggerFailure,
			})),
			Fn(`LET target = @target
RETURN WAITFOR EVENT "test" IN target
	TRIGGER (
		target <- "test"
	)
	TIMEOUT 20ms
	ON ERROR RETURN "error"`, expectTriggerObservable(waitFailure, "error", 1, 1, 1), "WAITFOR EVENT stream failure after trigger should clean up and use ON ERROR").Env(vm.WithParams(map[string]runtime.Value{
				"target": waitFailure,
			})),
			Fn(`LET target = @target
LET evt = WAITFOR EVENT "test" IN target
	TRIGGER (
		target <- "test"
	)
	TIMEOUT 20ms
	ON ERROR RETRY 2 DELAY 0 OR RETURN "error"
RETURN evt.type`, expectTriggerObservable(retry, "test", 2, 2, 2), "WAITFOR EVENT trigger should be retried through protected cleanup").Env(vm.WithParams(map[string]runtime.Value{
				"target": retry,
			})),
		}
	})
}

func expectTriggerObservable(target *TriggerObservable, expected any, subscribes, dispatches, closes int32) func(any) error {
	return func(actual any) error {
		if actual != expected {
			return fmt.Errorf("expected return value %v, got %v", expected, actual)
		}
		if got := target.SubscribeCount(); got != subscribes {
			return fmt.Errorf("expected %d subscribes, got %d", subscribes, got)
		}
		if got := target.DispatchCount(); got != dispatches {
			return fmt.Errorf("expected %d dispatches, got %d", dispatches, got)
		}
		if got := target.CloseCount(); got != closes {
			return fmt.Errorf("expected %d closes, got %d", closes, got)
		}

		return nil
	}
}
