package vm_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/compiler"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
	"github.com/MontFerret/ferret/v2/test/spec"
	. "github.com/MontFerret/ferret/v2/test/spec/exec"
)

func TestWaitforPredicate(t *testing.T) {
	RunSpecs(t, []spec.Spec{
		S(`
			LET start = NOW()
			LET ok = WAITFOR (DATE_DIFF(start, NOW(), "f") >= 30) TIMEOUT 0.5s EVERY 0.01s BACKOFF LINEAR
			RETURN ok
		`, true, "Should wait until predicate becomes true"),
		S(`
			LET start = NOW()
			LET ok = WAITFOR (DATE_DIFF(start, NOW(), "f") >= 30) TIMEOUT 0.5s EVERY 10ms, 30ms BACKOFF LINEAR JITTER 0.2
			RETURN ok
		`, true, "Should support EVERY cap with JITTER"),
		S(`
			LET ok = WAITFOR FALSE TIMEOUT 50ms EVERY 10ms
			RETURN ok
		`, false, "Should return false on timeout"),
		S(`
			LET ok = WAITFOR FALSE TIMEOUT 80ms EVERY 10ms, 10ms BACKOFF EXPONENTIAL JITTER 0.5
			RETURN ok
		`, false, "Should honor timeout with cap and jitter"),
		S(`
			LET start = NOW()
			LET ok = WAITFOR EXISTS (DATE_DIFF(start, NOW(), "f") > 20 ? [1] : []) TIMEOUT 0.5s EVERY 10ms
			RETURN ok
		`, true, "Should wait for non-empty value with EXISTS"),
		S(`
			LET start = NOW()
			LET ok = WAITFOR NOT EXISTS (DATE_DIFF(start, NOW(), "f") > 20 ? "" : "x") TIMEOUT 0.5s EVERY 10ms
			RETURN ok
		`, true, "Should wait for empty value with NOT EXISTS"),
		S(`
			LET start = NOW()
			LET token = WAITFOR VALUE (DATE_DIFF(start, NOW(), "f") > 20 ? "ok" : NONE) TIMEOUT 0.5s EVERY 10ms
			RETURN token
		`, "ok", "Should return value once it exists"),
		S(`
			LET start = NOW()
			LET ok = WAITFOR EXISTS (DATE_DIFF(start, NOW(), "f") > 20 ? { foo: 1 } : {}) TIMEOUT 0.5s EVERY 10ms
			RETURN ok
		`, true, "Should wait for non-empty object with EXISTS"),
		Object(`
			LET start = NOW()
			LET obj = WAITFOR VALUE (DATE_DIFF(start, NOW(), "f") > 20 ? { foo: 1 } : {}) TIMEOUT 0.5s EVERY 10ms
			RETURN obj
		`, map[string]any{"foo": 1}, "Should return object once it exists"),
		Array(`
			LET start = NOW()
			LET arr = WAITFOR VALUE (DATE_DIFF(start, NOW(), "f") > 20 ? [1, 2] : []) TIMEOUT 0.5s EVERY 10ms
			RETURN arr
		`, []any{1, 2}, "Should return array once it exists"),
		S(`
			LET start = NOW()
			LET ok = WAITFOR EXISTS (DATE_DIFF(start, NOW(), "f") > 20 ? "ok" : "") TIMEOUT 0.5s EVERY 10ms
			RETURN ok
		`, true, "Should wait for non-empty string with EXISTS"),
		S(`
			LET start = NOW()
			LET ok = WAITFOR NOT EXISTS (DATE_DIFF(start, NOW(), "f") > 20 ? {} : { foo: 1 }) TIMEOUT 0.5s EVERY 10ms
			RETURN ok
		`, true, "Should wait for empty object with NOT EXISTS"),
		S(`
			LET obj = {
				list: ["bar", "baz"],
			}
			LET ok = WAITFOR EXISTS obj.list TIMEOUT 0.5s EVERY 10ms
			RETURN ok
		`, true, "Should wait for non empty-array with EXISTS"),
		S(`
			LET obj = {
				list: [],
			}
			LET ok = WAITFOR NOT EXISTS obj.list TIMEOUT 0.5s EVERY 10ms
			RETURN ok
		`, true, "Should wait for non empty-array with EXISTS"),
		Nil(`
			LET token = WAITFOR VALUE NONE TIMEOUT 20ms EVERY 5ms
			RETURN token
		`, "Should return NONE on timeout for VALUE"),
		Nil(`
			LET token = WAITFOR VALUE NONE TIMEOUT 20ms EVERY 5ms ON ERROR FAIL
			RETURN token
		`, "Explicit FAIL should preserve timeout result semantics"),
		Object(`
			LET value = WAITFOR VALUE { ok: true, kind: "candidate" } WHEN .ok
			RETURN value
		`, map[string]any{"ok": true, "kind": "candidate"}, "WAITFOR VALUE WHEN should return the candidate value"),
		Object(`
			LET value = WAITFOR VALUE { ok: true, kind: "candidate" } WHEN .ok WHEN .kind == "candidate"
			RETURN value
		`, map[string]any{"ok": true, "kind": "candidate"}, "WAITFOR VALUE repeated WHEN clauses should return the candidate value"),
		S(`
			LET ok = WAITFOR EXISTS [1, 2, 3] WHEN LENGTH(.) >= 3
			RETURN ok
		`, true, "WAITFOR EXISTS WHEN should bind the full candidate array"),
		S(`
			LET ok = WAITFOR EXISTS [1, 2, 3] WHEN LENGTH(.) >= 3 WHEN LENGTH(.) > 0
			RETURN ok
		`, true, "WAITFOR EXISTS repeated WHEN clauses should bind the full candidate array"),
		S(`
			LET ok = WAITFOR NOT EXISTS [] WHEN LENGTH(.) == 0
			RETURN ok
		`, true, "WAITFOR NOT EXISTS WHEN should bind the not-existing candidate value"),
		S(`
			LET ok = WAITFOR NOT EXISTS [] WHEN LENGTH(.) == 0 WHEN . != NONE
			RETURN ok
		`, true, "WAITFOR NOT EXISTS repeated WHEN clauses should bind the not-existing candidate value"),
		Nil(`
			LET token = WAITFOR VALUE "ok" WHEN false TIMEOUT 20ms EVERY 5ms ON TIMEOUT RETURN NONE
			RETURN token
		`, "WAITFOR VALUE WHEN should honor ON TIMEOUT when the predicate never passes"),
		Nil(`
			LET token = WAITFOR VALUE "ok" WHEN true WHEN false TIMEOUT 20ms EVERY 5ms ON TIMEOUT RETURN NONE
			RETURN token
		`, "WAITFOR VALUE repeated WHEN clauses should honor ON TIMEOUT when any predicate never passes"),
		S(`
			LET ok = WAITFOR EXISTS [1] WHEN false TIMEOUT 20ms EVERY 5ms ON TIMEOUT RETURN false
			RETURN ok
		`, false, "WAITFOR EXISTS WHEN should honor ON TIMEOUT when the predicate never passes"),
	})
}

func TestWaitforPredicateWhenRetriesUntilTrue(t *testing.T) {
	for _, level := range []compiler.OptimizationLevel{compiler.O0, compiler.O1} {
		callCount := 0

		RunSpecsWith(
			t,
			fmt.Sprintf("VM/O%d", level),
			compiler.New(compiler.WithOptimizationLevel(level)),
			[]spec.Spec{
				S(`
					LET token = WAITFOR VALUE CANDIDATE() WHEN .state == "ready" TIMEOUT 100ms EVERY 0
					RETURN token.value
				`, "ok", "WAITFOR VALUE WHEN should retry until the predicate passes"),
			},
			vm.WithFunction("CANDIDATE", func(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
				callCount++

				state := "pending"
				if callCount >= 3 {
					state = "ready"
				}

				return runtime.NewObjectWith(map[string]runtime.Value{
					"state": runtime.NewString(state),
					"value": runtime.NewString("ok"),
				}), nil
			}),
		)

		if got, want := callCount, 3; got != want {
			t.Fatalf("unexpected WAITFOR VALUE WHEN candidate call count for O%d: got %d, want %d", level, got, want)
		}
	}
}

func TestWaitforPredicateWhenSkipsPredicateUntilBasePasses(t *testing.T) {
	for _, level := range []compiler.OptimizationLevel{compiler.O0, compiler.O1} {
		predicateCalls := 0

		RunSpecsWith(
			t,
			fmt.Sprintf("VM/O%d", level),
			compiler.New(compiler.WithOptimizationLevel(level)),
			[]spec.Spec{
				S(`
					LET ok = WAITFOR EXISTS NONE WHEN PREDICATE(.) TIMEOUT 20ms EVERY 1ms ON TIMEOUT RETURN false
					RETURN ok
				`, false, "WAITFOR EXISTS WHEN should not evaluate the predicate before existence passes"),
			},
			vm.WithFunction("PREDICATE", func(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
				predicateCalls++

				return runtime.True, nil
			}),
		)

		if got := predicateCalls; got != 0 {
			t.Fatalf("WAITFOR EXISTS WHEN evaluated predicate before existence passed for O%d: got %d calls", level, got)
		}
	}
}

func TestWaitforPredicateMultipleWhenShortCircuits(t *testing.T) {
	for _, level := range []compiler.OptimizationLevel{compiler.O0, compiler.O1} {
		firstCalls := 0
		secondCalls := 0

		RunSpecsWith(
			t,
			fmt.Sprintf("VM/O%d", level),
			compiler.New(compiler.WithOptimizationLevel(level)),
			[]spec.Spec{
				Nil(`
					LET token = WAITFOR VALUE "ok" WHEN REPEATED_WHEN_FIRST(.) WHEN REPEATED_WHEN_SECOND(.) TIMEOUT 20ms EVERY 1ms ON TIMEOUT RETURN NONE
					RETURN token
				`, "WAITFOR VALUE repeated WHEN clauses should short-circuit after the first false predicate"),
			},
			vm.WithFunction("REPEATED_WHEN_FIRST", func(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
				firstCalls++

				return runtime.False, nil
			}),
			vm.WithFunction("REPEATED_WHEN_SECOND", func(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
				secondCalls++

				return runtime.True, nil
			}),
		)

		if firstCalls == 0 {
			t.Fatalf("WAITFOR VALUE repeated WHEN did not evaluate the first predicate for O%d", level)
		}
		if secondCalls != 0 {
			t.Fatalf("WAITFOR VALUE repeated WHEN evaluated a later predicate after false for O%d: got %d calls", level, secondCalls)
		}
	}
}

func TestWaitforPredicateWhenUsesOperationErrorPolicy(t *testing.T) {
	RunSpecs(t, []spec.Spec{
		S(`
			RETURN WAITFOR VALUE "ok" WHEN FAIL_PREDICATE(.) TIMEOUT 20ms EVERY 0 ON ERROR RETURN "error"
		`, "error", "WAITFOR VALUE WHEN predicate errors should use the wait error policy"),
		S(`
			RETURN WAITFOR VALUE "ok" WHEN true WHEN FAIL_PREDICATE(.) TIMEOUT 20ms EVERY 0 ON ERROR RETURN "error"
		`, "error", "WAITFOR VALUE repeated WHEN predicate errors should use the wait error policy"),
	}, vm.WithFunction("FAIL_PREDICATE", func(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
		return runtime.None, fmt.Errorf("predicate failed")
	}))
}
