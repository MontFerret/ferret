package vm_test

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/MontFerret/ferret/v2/pkg/compiler"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/source"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

type cancellationStringValue struct {
	calls *int
}

func (v *cancellationStringValue) String() string {
	(*v.calls)++
	return "probe"
}

func (*cancellationStringValue) Hash() uint64 {
	return 1
}

func (v *cancellationStringValue) Copy() runtime.Value {
	return v
}

func (*cancellationStringValue) Type() runtime.Type {
	return runtime.TypeObject
}

func TestPureExecutionCancellationAtBackwardJumpSafepoint(t *testing.T) {
	const query = `
LET signal = START()
VAR i = 0
FOR WHILE i < 1000000000
  i = i + 1
RETURN signal
`

	for _, level := range []compiler.OptimizationLevel{compiler.O0, compiler.O1} {
		t.Run(optimizationName(level), func(t *testing.T) {
			instance := compileCancellationVM(t, level, query)
			started := make(chan struct{})
			var startedOnce sync.Once
			env := cancellationEnvironment(t, vm.WithFunction("START", func(context.Context, ...runtime.Value) (runtime.Value, error) {
				startedOnce.Do(func() { close(started) })
				return runtime.True, nil
			}))

			ctx, cancel := context.WithCancel(context.Background())
			done := make(chan error, 1)
			go func() {
				result, err := instance.Run(ctx, env)
				if result != nil {
					_ = result.Close()
				}
				done <- err
			}()

			select {
			case <-started:
				cancel()
			case <-time.After(5 * time.Second):
				cancel()
				t.Fatal("execution did not enter the loop")
			}

			select {
			case err := <-done:
				if !errors.Is(err, context.Canceled) {
					t.Fatalf("run error = %v, want context.Canceled", err)
				}
			case <-time.After(5 * time.Second):
				t.Fatal("execution did not stop after cancellation")
			}
		})
	}
}

func TestReturnedUDFForCancellationAtBackwardJumpSafepoint(t *testing.T) {
	const query = `
FUNC spin() {
  LET signal = START()
  VAR i = 0
  RETURN FOR WHILE i < 1000000000 {
    i = i + 1
    RETURN signal
  }
}
RETURN spin()
`

	for _, level := range []compiler.OptimizationLevel{compiler.O0, compiler.O1} {
		t.Run(optimizationName(level), func(t *testing.T) {
			instance := compileCancellationVM(t, level, query)
			started := make(chan struct{})
			var startedOnce sync.Once
			env := cancellationEnvironment(t, vm.WithFunction("START", func(context.Context, ...runtime.Value) (runtime.Value, error) {
				startedOnce.Do(func() { close(started) })
				return runtime.True, nil
			}))

			ctx, cancel := context.WithCancel(context.Background())
			done := make(chan error, 1)
			go func() {
				result, err := instance.Run(ctx, env)
				if result != nil {
					_ = result.Close()
				}
				done <- err
			}()

			select {
			case <-started:
				cancel()
			case <-time.After(5 * time.Second):
				cancel()
				t.Fatal("execution did not enter the UDF loop")
			}

			select {
			case err := <-done:
				if !errors.Is(err, context.Canceled) {
					t.Fatalf("run error = %v, want context.Canceled", err)
				}
			case <-time.After(5 * time.Second):
				t.Fatal("execution did not stop after cancellation")
			}
		})
	}
}

func TestPreCanceledExecutionStopsAtCallAndReturnSafepoints(t *testing.T) {
	tests := map[string]string{
		"return":   `RETURN 1 + 2`,
		"udf call": "FUNC value() => 1\nRETURN value()",
	}

	for _, level := range []compiler.OptimizationLevel{compiler.O0, compiler.O1} {
		for name, query := range tests {
			t.Run(optimizationName(level)+"/"+name, func(t *testing.T) {
				instance := compileCancellationVM(t, level, query)
				ctx, cancel := context.WithCancel(context.Background())
				cancel()

				result, err := instance.Run(ctx, cancellationEnvironment(t))
				if result != nil {
					_ = result.Close()
				}
				if !errors.Is(err, context.Canceled) {
					t.Fatalf("run error = %v, want context.Canceled", err)
				}
			})
		}
	}
}

func TestNestedReturnDoesNotPollCancellation(t *testing.T) {
	const query = `
FUNC inner() => CANCEL()
LET value = inner()
RETURN "after-" + @probe
`

	for _, level := range []compiler.OptimizationLevel{compiler.O0, compiler.O1} {
		t.Run(optimizationName(level), func(t *testing.T) {
			instance := compileCancellationVM(t, level, query)
			ctx, cancel := context.WithCancel(context.Background())
			stringCalls := 0
			probe := &cancellationStringValue{calls: &stringCalls}
			env := cancellationEnvironment(t,
				vm.WithParam("probe", probe),
				vm.WithFunction("CANCEL", func(context.Context, ...runtime.Value) (runtime.Value, error) {
					cancel()
					return runtime.None, nil
				}),
			)

			result, err := instance.Run(ctx, env)
			if result != nil {
				_ = result.Close()
			}
			if !errors.Is(err, context.Canceled) {
				t.Fatalf("run error = %v, want context.Canceled", err)
			}
			if stringCalls != 1 {
				t.Fatalf("post-return string calls = %d, want 1", stringCalls)
			}
		})
	}
}

func TestCancellationBypassesProtectedRecovery(t *testing.T) {
	const query = `RETURN CANCEL() ON ERROR RETURN AFTER()`
	tests := []struct {
		err       error
		name      string
		cancelCtx bool
	}{
		{name: "canceled context", err: context.Canceled, cancelCtx: true},
		{name: "returned cancellation", err: fmt.Errorf("host canceled: %w", context.Canceled)},
		{name: "returned deadline", err: fmt.Errorf("host deadline: %w", context.DeadlineExceeded)},
	}

	for _, level := range []compiler.OptimizationLevel{compiler.O0, compiler.O1} {
		for _, test := range tests {
			t.Run(optimizationName(level)+"/"+test.name, func(t *testing.T) {
				instance := compileCancellationVM(t, level, query)
				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()
				var afterCalls atomic.Int64
				env := cancellationEnvironment(t,
					vm.WithFunction("CANCEL", func(context.Context, ...runtime.Value) (runtime.Value, error) {
						if test.cancelCtx {
							cancel()
						}

						return runtime.None, test.err
					}),
					vm.WithFunction("AFTER", func(context.Context, ...runtime.Value) (runtime.Value, error) {
						afterCalls.Add(1)
						return runtime.NewInt(42), nil
					}),
				)

				result, err := instance.Run(ctx, env)
				if result != nil {
					_ = result.Close()
				}
				if !errors.Is(err, test.err) {
					t.Fatalf("run error = %v, want %v", err, test.err)
				}
				if calls := afterCalls.Load(); calls != 0 {
					t.Fatalf("recovery calls = %d, want zero", calls)
				}
			})
		}
	}
}

func compileCancellationVM(t *testing.T, level compiler.OptimizationLevel, query string) *vm.VM {
	t.Helper()

	program, err := compiler.New(compiler.WithOptimizationLevel(level)).Compile(source.NewAnonymous(query))
	if err != nil {
		t.Fatalf("compile query: %v", err)
	}
	instance, err := vm.New(program)
	if err != nil {
		t.Fatalf("create VM: %v", err)
	}
	t.Cleanup(func() { _ = instance.Close() })

	return instance
}

func cancellationEnvironment(t *testing.T, options ...vm.EnvironmentOption) *vm.Environment {
	t.Helper()

	env, err := vm.NewEnvironment(options)
	if err != nil {
		t.Fatalf("create environment: %v", err)
	}

	return env
}

func optimizationName(level compiler.OptimizationLevel) string {
	if level == compiler.O0 {
		return "O0"
	}

	return "O1"
}
