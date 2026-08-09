package vm_test

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

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
LET base = 0.5ms

LET status = WAITFOR EVENT "test" IN obs TIMEOUT base * 2 ON TIMEOUT RETURN "timeout"

RETURN status`, "timeout", "WAITFOR EVENT should accept a computed native duration timeout").Env(vm.WithParams(map[string]runtime.Value{
			"obs": blocking,
		})),
		S(`LET obs = @obs
RETURN WAITFOR EVENT "test" IN obs TIMEOUT 1 ON TIMEOUT RETURN "timeout"`, "timeout", "WAITFOR EVENT should treat numeric timeouts as milliseconds").Env(vm.WithParams(map[string]runtime.Value{
			"obs": blocking,
		})),
		Error(`LET obs = @obs
LET timeout = -1ms
RETURN WAITFOR EVENT "test" IN obs TIMEOUT timeout`, "WAITFOR EVENT should reject negative timeouts").Env(vm.WithParams(map[string]runtime.Value{
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

func TestWaitforEventSynchronizationGroups(t *testing.T) {
	anyFirst := NewObservable([]runtime.Value{NewTestEventType("first")})
	anySecond := NewObservable([]runtime.Value{NewTestEventType("second")})
	filtered := NewObservable([]runtime.Value{
		NewTestEventType("ignored"),
		NewTestEventType("accepted"),
	})
	allFirst := NewObservable([]runtime.Value{NewTestEventType("first")})
	allSecond := NewObservable([]runtime.Value{NewTestEventType("second")})
	blocking := NewBlockingObservable()

	RunSpecs(t, []spec.Spec{
		Fn(`LET first = @first
LET second = @second
LET evt = WAITFOR EVENT ANY {
	"first" IN first
	"second" IN second
}
RETURN evt.type`, func(actual any) error {
			if actual != "first" && actual != "second" {
				return fmt.Errorf("expected either event winner, got %v", actual)
			}
			return nil
		}, "WAITFOR EVENT ANY should return the naturally selected winner").Env(vm.WithParams(map[string]runtime.Value{
			"first": anyFirst, "second": anySecond,
		})),
		S(`LET filtered = @filtered
LET blocking = @blocking
LET evt = WAITFOR EVENT ANY {
	"filtered" IN filtered WHEN .type == "accepted"
	"blocking" IN blocking
}
RETURN evt.type`, "accepted", "WAITFOR EVENT ANY should apply filters only to the yielding arm").Env(vm.WithParams(map[string]runtime.Value{
			"filtered": filtered, "blocking": blocking,
		})),
		Array(`LET first = @first
LET second = @second
RETURN WAITFOR EVENT ALL {
	"first" IN first
	"second" IN second
}`, []any{
			map[string]any{"type": "first"},
			map[string]any{"type": "second"},
		}, "WAITFOR EVENT ALL should return declaration order").Env(vm.WithParams(map[string]runtime.Value{
			"first": allFirst, "second": allSecond,
		})),
		S(`LET first = @first
LET blocking = @blocking
RETURN WAITFOR EVENT ALL {
	"first" IN first
	"blocking" IN blocking
} TIMEOUT 2ms ON TIMEOUT RETURN "timeout"`, "timeout", "WAITFOR EVENT ALL should use one shared deadline").Env(vm.WithParams(map[string]runtime.Value{
			"first": allFirst, "blocking": blocking,
		})),
	})
}

func TestWaitforEventSynchronizationSameSource(t *testing.T) {
	RunSpecFactory(t, func() []spec.Spec {
		target := NewTriggerObservable()

		return []spec.Spec{
			Fn(`LET target = @target
LET evt = WAITFOR EVENT ANY {
	"first" IN target
	"second" IN target
} TRIGGER target <- "first" TIMEOUT 20ms
RETURN evt.type`, expectTriggerObservable(target, "first", 2, 1, 2), "WAITFOR EVENT ANY should establish and clean up every same-source subscription").Env(vm.WithParams(map[string]runtime.Value{
				"target": target,
			})),
		}
	})
}

func TestWaitforEventSynchronizationTrigger(t *testing.T) {
	RunSpecFactory(t, func() []spec.Spec {
		anyFirst := NewTriggerObservable()
		anySecond := NewTriggerObservable()
		allFirst := NewTriggerObservable()
		allSecond := NewTriggerObservable()
		triggerFailureFirst := NewTriggerObservable()
		triggerFailureSecond := NewTriggerObservable()
		triggerFailureFirst.FailNextDispatches(1, errors.New("group trigger failed"))
		readFailureFirst := NewTriggerObservable()
		readFailureSecond := NewTriggerObservable()
		readFailureFirst.FailReadsWith(errors.New("group stream failed"))
		retryFirst := NewTriggerObservable()
		retrySecond := NewTriggerObservable()
		retryFirst.FailNextDispatches(1, errors.New("group trigger failed once"))

		return []spec.Spec{
			Fn(`LET first = @first
LET second = @second
LET evt = WAITFOR EVENT ANY {
	"first" IN first
	"second" IN second
} TRIGGER (
	first <- "first"
	second <- "second"
) TIMEOUT 20ms
RETURN evt.type`, expectTriggerGroup(anyFirst, anySecond, []any{"first", "second"}, 1, 1), "WAITFOR EVENT ANY should arm every subscription before one trigger").Env(vm.WithParams(map[string]runtime.Value{
				"first": anyFirst, "second": anySecond,
			})),
			Fn(`LET first = @first
LET second = @second
RETURN WAITFOR EVENT ALL {
	"first" IN first
	"second" IN second
} TRIGGER (
	second <- "second"
	first <- "first"
) TIMEOUT 20ms`, expectTriggerGroup(allFirst, allSecond, []any{
				map[string]any{"type": "first"},
				map[string]any{"type": "second"},
			}, 1, 1), "WAITFOR EVENT ALL should collect occurrence order but return declaration order").Env(vm.WithParams(map[string]runtime.Value{
				"first": allFirst, "second": allSecond,
			})),
			Fn(`LET first = @first
LET second = @second
RETURN WAITFOR EVENT ANY {
	"first" IN first
	"second" IN second
} TRIGGER (
	first <- "first"
	second <- "second"
) TIMEOUT 20ms ON ERROR RETURN "error"`, expectTriggerGroup(triggerFailureFirst, triggerFailureSecond, "error", 1, 1), "WAITFOR EVENT group trigger failure should close every subscription").Env(vm.WithParams(map[string]runtime.Value{
				"first": triggerFailureFirst, "second": triggerFailureSecond,
			})),
			Fn(`LET first = @first
LET second = @second
RETURN WAITFOR EVENT ANY {
	"first" IN first
	"second" IN second
} TRIGGER (
	first <- "first"
) TIMEOUT 20ms ON ERROR RETURN "error"`, expectTriggerGroup(readFailureFirst, readFailureSecond, "error", 1, 1), "active event-group arm errors should fail and close the group").Env(vm.WithParams(map[string]runtime.Value{
				"first": readFailureFirst, "second": readFailureSecond,
			})),
			Fn(`LET first = @first
LET second = @second
LET evt = WAITFOR EVENT ANY {
	"first" IN first
	"second" IN second
} TRIGGER (
	first <- "first"
	second <- "second"
) TIMEOUT 20ms ON ERROR RETRY 2 DELAY 0s OR RETURN "error"
RETURN evt.type`, expectTriggerGroup(retryFirst, retrySecond, []any{"first", "second"}, 2, 2), "WAITFOR EVENT groups should retry after full cleanup").Env(vm.WithParams(map[string]runtime.Value{
				"first": retryFirst, "second": retrySecond,
			})),
		}
	})
}

func TestWaitforEventSynchronizationCancellationClosesAll(t *testing.T) {
	first := NewTriggerObservable()
	second := NewTriggerObservable()
	program, err := spec.Compile(`LET first = @first
LET second = @second
RETURN WAITFOR EVENT ALL {
	"first" IN first
	"second" IN second
}`)
	if err != nil {
		t.Fatalf("compile event group: %v", err)
	}

	instance, err := vm.NewWith(program)
	if err != nil {
		t.Fatalf("create VM: %v", err)
	}
	t.Cleanup(func() { _ = instance.Close() })
	environment, err := vm.NewEnvironment([]vm.EnvironmentOption{
		vm.WithNamespace(spec.Stdlib()),
		vm.WithParams(map[string]runtime.Value{"first": first, "second": second}),
	})
	if err != nil {
		t.Fatalf("create environment: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() {
		result, runErr := instance.Run(ctx, environment)
		if result != nil {
			_ = result.Close()
		}
		done <- runErr
	}()

	deadline := time.After(time.Second)
	for first.SubscribeCount() != 1 || second.SubscribeCount() != 1 {
		select {
		case <-deadline:
			t.Fatal("subscriptions were not established")
		case <-time.After(time.Millisecond):
		}
	}
	cancel()

	select {
	case runErr := <-done:
		if !errors.Is(runErr, context.Canceled) {
			t.Fatalf("expected cancellation, got %v", runErr)
		}
	case <-time.After(time.Second):
		t.Fatal("event group did not stop after cancellation")
	}

	assertTriggerObservableCounts(t, first, 1, 0, 1)
	assertTriggerObservableCounts(t, second, 1, 0, 1)
}

func TestWaitforEventTrigger(t *testing.T) {
	RunSpecFactory(t, func() []spec.Spec {
		armed := NewTriggerObservable()
		timeout := NewTriggerObservable()
		triggerFailure := NewTriggerObservable()
		triggerFailure.FailNextDispatches(1, errors.New("trigger failed"))
		triggerCallFailure := NewTriggerObservable()
		waitFailure := NewTriggerObservable()
		waitFailure.FailReadsWith(errors.New("stream failed"))
		retry := NewTriggerObservable()
		retry.FailNextDispatches(1, errors.New("trigger failed once"))

		return []spec.Spec{
			Fn(`LET target = @target
LET evt = WAITFOR EVENT "test" IN target
	TRIGGER target <- "test"
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
	TRIGGER FAIL()
	TIMEOUT 1ms
	ON ERROR RETURN "error"`, expectTriggerObservable(triggerCallFailure, "error", 1, 0, 1), "WAITFOR EVENT inline trigger call failure should belong to outer recovery").Env(
				vm.WithParams(map[string]runtime.Value{
					"target": triggerCallFailure,
				}),
				vm.WithFunction("FAIL", func(context.Context, ...runtime.Value) (runtime.Value, error) {
					return runtime.None, errors.New("trigger failed")
				}),
			),
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
	ON ERROR RETRY 2 DELAY 0s OR RETURN "error"
RETURN evt.type`, expectTriggerObservable(retry, "test", 2, 2, 2), "WAITFOR EVENT trigger should be retried through protected cleanup").Env(vm.WithParams(map[string]runtime.Value{
				"target": retry,
			})),
		}
	})
}

func TestWaitforEventTriggerCleanupOnTriggerError(t *testing.T) {
	failFn := vm.WithFunction("FAIL", func(context.Context, ...runtime.Value) (runtime.Value, error) {
		return runtime.None, errors.New("trigger failed")
	})

	t.Run("plain trigger dispatch failure closes before returning error", func(t *testing.T) {
		target := NewTriggerObservable()
		target.FailNextDispatches(1, errors.New("trigger failed"))

		result, err := runWaitforEventTriggerProgram(t, `LET target = @target
RETURN WAITFOR EVENT "test" IN target
	TRIGGER target <- "test"`, target)
		if result != nil {
			_ = result.Close()
		}
		assertTriggerRuntimeError(t, err, "trigger failed")
		assertTriggerObservableCounts(t, target, 1, 1, 1)
	})

	t.Run("plain trigger call failure closes before returning error", func(t *testing.T) {
		target := NewTriggerObservable()

		result, err := runWaitforEventTriggerProgram(t, `LET target = @target
RETURN WAITFOR EVENT "test" IN target
	TRIGGER FAIL()`, target, failFn)
		if result != nil {
			_ = result.Close()
		}
		assertTriggerRuntimeError(t, err, "trigger failed")
		assertTriggerObservableCounts(t, target, 1, 0, 1)
	})

	t.Run("timeout-only trigger call failure closes before returning error", func(t *testing.T) {
		target := NewTriggerObservable()

		result, err := runWaitforEventTriggerProgram(t, `LET target = @target
RETURN WAITFOR EVENT "test" IN target
	TRIGGER FAIL()
	TIMEOUT 1ms
	ON TIMEOUT RETURN "timeout"`, target, failFn)
		if result != nil {
			_ = result.Close()
		}
		assertTriggerRuntimeError(t, err, "trigger failed")
		assertTriggerObservableCounts(t, target, 1, 0, 1)
	})

	t.Run("outer suppression closes before result close", func(t *testing.T) {
		target := NewTriggerObservable()

		result, err := runWaitforEventTriggerProgram(t, `LET target = @target
LET out = (WAITFOR EVENT "test" IN target
	TRIGGER FAIL())?
RETURN out`, target, failFn)
		if err != nil {
			t.Fatalf("expected suppressed trigger failure, got %v", err)
		}
		defer func() {
			_ = result.Close()
		}()

		assertTriggerObservableCounts(t, target, 1, 0, 1)
	})

	t.Run("timeout-aware outer suppression closes before result close", func(t *testing.T) {
		target := NewTriggerObservable()

		result, err := runWaitforEventTriggerProgram(t, `LET target = @target
LET out = (WAITFOR EVENT "test" IN target
	TRIGGER FAIL()
	TIMEOUT 1ms
	ON TIMEOUT RETURN "timeout")?
RETURN out`, target, failFn)
		if err != nil {
			t.Fatalf("expected suppressed trigger failure, got %v", err)
		}
		defer func() {
			_ = result.Close()
		}()

		assertTriggerObservableCounts(t, target, 1, 0, 1)
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

func expectTriggerGroup(first, second *TriggerObservable, expected any, subscribes, closes int32) func(any) error {
	return func(actual any) error {
		if expectedValues, ok := expected.([]any); ok && len(expectedValues) == 2 {
			if actualValues, arrayOK := actual.([]any); arrayOK {
				if fmt.Sprint(actualValues) != fmt.Sprint(expectedValues) {
					return fmt.Errorf("expected return value %v, got %v", expected, actual)
				}
			} else if actual != expectedValues[0] && actual != expectedValues[1] {
				return fmt.Errorf("expected return value in %v, got %v", expectedValues, actual)
			}
		} else if actual != expected {
			return fmt.Errorf("expected return value %v, got %v", expected, actual)
		}

		for idx, target := range []*TriggerObservable{first, second} {
			if got := target.SubscribeCount(); got != subscribes {
				return fmt.Errorf("target %d: expected %d subscribes, got %d", idx, subscribes, got)
			}
			if got := target.CloseCount(); got != closes {
				return fmt.Errorf("target %d: expected %d closes, got %d", idx, closes, got)
			}
		}

		return nil
	}
}

func runWaitforEventTriggerProgram(
	t *testing.T,
	query string,
	target *TriggerObservable,
	opts ...vm.EnvironmentOption,
) (*vm.Result, error) {
	t.Helper()

	prog, err := spec.Compile(query)
	if err != nil {
		t.Fatalf("compile failed: %v", err)
	}

	instance, err := vm.NewWith(prog)
	if err != nil {
		t.Fatalf("vm init failed: %v", err)
	}
	t.Cleanup(func() {
		_ = instance.Close()
	})

	envOpts := []vm.EnvironmentOption{
		vm.WithNamespace(spec.Stdlib()),
		vm.WithParam("target", target),
	}
	envOpts = append(envOpts, opts...)

	env, err := vm.NewEnvironment(envOpts)
	if err != nil {
		t.Fatalf("environment build failed: %v", err)
	}

	return instance.Run(context.Background(), env)
}

func assertTriggerRuntimeError(t *testing.T, err error, contains string) {
	t.Helper()

	if err == nil {
		t.Fatal("expected runtime error")
	}

	var rtErr *vm.RuntimeError
	if !errors.As(err, &rtErr) {
		t.Fatalf("expected runtime error, got %T", err)
	}

	if !strings.Contains(rtErr.Format(), contains) {
		t.Fatalf("expected runtime error to contain %q, got:\n%s", contains, rtErr.Format())
	}
}

func assertTriggerObservableCounts(t *testing.T, target *TriggerObservable, subscribes, dispatches, closes int32) {
	t.Helper()

	if got := target.SubscribeCount(); got != subscribes {
		t.Fatalf("expected %d subscribes, got %d", subscribes, got)
	}
	if got := target.DispatchCount(); got != dispatches {
		t.Fatalf("expected %d dispatches, got %d", dispatches, got)
	}
	if got := target.CloseCount(); got != closes {
		t.Fatalf("expected %d closes, got %d", closes, got)
	}
}
