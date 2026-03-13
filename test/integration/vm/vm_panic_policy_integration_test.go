package vm_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

func TestPanicPolicyRecoversPanics(t *testing.T) {
	forEachCompiledVM(t, "RETURN PANIC_FN()", func(t *testing.T, instance *vm.VM) {
		env := mustNewEnvironment(t, vm.WithFunction("PANIC_FN", func(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
			panic("panic in host function")
		}))

		_, err := runVM(t, instance, env)
		if err == nil {
			t.Fatal("expected runtime error")
		}

		var rtErr *vm.RuntimeError
		if !errors.As(err, &rtErr) {
			t.Fatalf("expected runtime error, got %T", err)
		}
	})
}

func TestPanicPolicyPropagatesPanics(t *testing.T) {
	forEachCompiledVMWithOptions(
		t,
		"RETURN PANIC_FN()",
		[]vm.Option{vm.WithPanicPolicy(vm.PanicPropagate)},
		func(t *testing.T, instance *vm.VM) {
			env := mustNewEnvironment(t, vm.WithFunction("PANIC_FN", func(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
				panic("panic in host function")
			}))

			defer func() {
				if recovered := recover(); recovered == nil {
					t.Fatal("expected panic to propagate")
				}
			}()

			_, _ = runVM(t, instance, env)
		},
	)
}

func TestPanicPolicyPropagateStillWrapsReturnedErrors(t *testing.T) {
	forEachCompiledVMWithOptions(
		t,
		"RETURN FAIL_FN()",
		[]vm.Option{vm.WithPanicPolicy(vm.PanicPropagate)},
		func(t *testing.T, instance *vm.VM) {
			env := mustNewEnvironment(t, vm.WithFunction("FAIL_FN", func(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
				return runtime.None, errors.New("boom")
			}))

			_, err := runVM(t, instance, env)
			if err == nil {
				t.Fatal("expected runtime error")
			}

			var rtErr *vm.RuntimeError
			if !errors.As(err, &rtErr) {
				t.Fatalf("expected runtime error, got %T", err)
			}
		},
	)
}

func TestRecoveredPanicRuntimeErrorDoesNotLeakGoStackTrace(t *testing.T) {
	forEachCompiledVM(t, "RETURN PANIC_FN()", func(t *testing.T, instance *vm.VM) {
		env := mustNewEnvironment(t, vm.WithFunction("PANIC_FN", func(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
			panic("panic in host function")
		}))

		_, err := runVM(t, instance, env)
		if err == nil {
			t.Fatal("expected runtime error")
		}

		var rtErr *vm.RuntimeError
		if !errors.As(err, &rtErr) {
			t.Fatalf("expected runtime error, got %T", err)
		}

		formatted := rtErr.Format()
		if strings.Contains(formatted, "goroutine ") {
			t.Fatalf("unexpected goroutine stack trace leak:\n%s", formatted)
		}

		if strings.Contains(formatted, "runtime/panic.go") {
			t.Fatalf("unexpected Go runtime stack trace leak:\n%s", formatted)
		}
	})
}
