package vm_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
	"github.com/MontFerret/ferret/v2/test/spec"
	"github.com/MontFerret/ferret/v2/test/spec/assert"
	. "github.com/MontFerret/ferret/v2/test/spec/exec"
)

func TestPanicPolicyRecoversPanics(t *testing.T) {
	RunSpecFactory(t, func() []spec.Spec {
		return []spec.Spec{
			spec.NewSpec("RETURN PANIC_FN()").
				Env(vm.WithFunction("PANIC_FN", func(context.Context, ...runtime.Value) (runtime.Value, error) {
					panic("panic in host function")
				})).
				Expect().ExecError(ShouldBeRuntimeError, &ExpectedRuntimeError{}),
		}
	})
}

func TestPanicPolicyPropagatesPanics(t *testing.T) {
	RunSequenceFactory(t, func() []spec.Sequence {
		return []spec.Sequence{
			{
				Base: spec.NewBaseSpec("RETURN PANIC_FN()"),
				VM: []vm.Option{
					vm.WithPanicPolicy(vm.PanicPropagate),
				},
				Steps: []spec.SequenceStep{
					{
						Name: "panic propagates",
						Env: []vm.EnvironmentOption{
							vm.WithFunction("PANIC_FN", func(context.Context, ...runtime.Value) (runtime.Value, error) {
								panic("panic in host function")
							}),
						},
						Panic: spec.NewExpectation(assert.ShouldEqual, "panic in host function"),
					},
				},
			},
		}
	})
}

func TestPanicPolicyPropagateStillWrapsReturnedErrors(t *testing.T) {
	RunSpecFactory(t, func() []spec.Spec {
		return []spec.Spec{
			spec.NewSpec("RETURN FAIL_FN()").
				VM(vm.WithPanicPolicy(vm.PanicPropagate)).
				Env(vm.WithFunction("FAIL_FN", func(context.Context, ...runtime.Value) (runtime.Value, error) {
					return runtime.None, errors.New("boom")
				})).
				Expect().ExecError(ShouldBeRuntimeError, &ExpectedRuntimeError{}),
		}
	})
}

func TestRecoveredPanicRuntimeErrorDoesNotLeakGoStackTrace(t *testing.T) {
	RunSpecFactory(t, func() []spec.Spec {
		return []spec.Spec{
			spec.NewSpec("RETURN PANIC_FN()").
				Env(vm.WithFunction("PANIC_FN", func(context.Context, ...runtime.Value) (runtime.Value, error) {
					panic("panic in host function")
				})).
				Expect().ExecError(ShouldBeRuntimeError, &ExpectedRuntimeError{
				NotContains: []string{"goroutine ", "runtime/panic.go"},
			}),
		}
	})
}

func TestStreamGroupCallbackPanicsFollowPanicPolicy(t *testing.T) {
	tests := []struct {
		name       string
		query      string
		params     func() (map[string]runtime.Value, *panicPolicyObservable)
		panicValue string
	}{
		{
			name: "subscribe",
			query: `LET ready = @ready
LET panicSource = @panic
RETURN WAITFOR EVENT ANY {
	"ready" IN ready
	"panic" IN panicSource
}`,
			params: func() (map[string]runtime.Value, *panicPolicyObservable) {
				ready := &panicPolicyObservable{}
				return map[string]runtime.Value{
					"ready": ready,
					"panic": &panicPolicyObservable{subscribePanic: "subscribe panic"},
				}, ready
			},
			panicValue: "subscribe panic",
		},
		{
			name: "read",
			query: `LET panicSource = @panic
RETURN WAITFOR EVENT ANY { "panic" IN panicSource }`,
			panicValue: "read panic",
			params: func() (map[string]runtime.Value, *panicPolicyObservable) {
				observable := &panicPolicyObservable{readPanic: "read panic"}
				return map[string]runtime.Value{"panic": observable}, observable
			},
		},
		{
			name: "partial close",
			query: `LET closeSource = @close
LET failure = @failure
RETURN WAITFOR EVENT ANY {
	"close" IN closeSource
	"failure" IN failure
}`,
			params: func() (map[string]runtime.Value, *panicPolicyObservable) {
				closer := &panicPolicyObservable{closePanic: "close panic"}
				return map[string]runtime.Value{
					"close":   closer,
					"failure": &panicPolicyObservable{subscribeErr: errors.New("setup failed")},
				}, closer
			},
			panicValue: "close panic",
		},
	}

	for _, test := range tests {
		for _, policy := range []struct {
			name string
			mode vm.PanicPolicy
		}{
			{name: "recover", mode: vm.PanicRecover},
			{name: "propagate", mode: vm.PanicPropagate},
		} {
			t.Run(test.name+"/"+policy.name, func(t *testing.T) {
				params, closed := test.params()
				recovered, runErr := runStreamGroupPanic(t, test.query, params, policy.mode)
				if policy.mode == vm.PanicRecover {
					if runErr == nil {
						t.Fatalf("expected recovered panic %q, got %v (close count %d)", test.panicValue, runErr, closed.closeCount.Load())
					}

					if !strings.Contains(runErr.Error(), "unexpected runtime panic") {
						t.Fatalf("expected runtime panic error, got %v", runErr)
					}

					if recovered != nil {
						t.Fatalf("recovered policy propagated panic %v", recovered)
					}
				} else {
					if recovered != test.panicValue {
						t.Fatalf("expected propagated panic %q, got %v", test.panicValue, recovered)
					}
				}

				if got := closed.closeCount.Load(); got != 1 {
					t.Fatalf("expected established stream close once, got %d", got)
				}
			})
		}
	}
}

func TestStreamGroupInternalSetupCancellationUsesOnError(t *testing.T) {
	setupErr := errors.New("setup failed")
	params := map[string]runtime.Value{
		"failure": &panicPolicyObservable{subscribeErr: setupErr},
		"peer":    &panicPolicyObservable{waitForCancel: true},
	}
	program, err := spec.Compile(`LET failure = @failure
LET peer = @peer
RETURN WAITFOR EVENT ANY {
	"failure" IN failure
	"peer" IN peer
} ON ERROR RETURN "recovered"`)
	if err != nil {
		t.Fatalf("compile grouped wait: %v", err)
	}

	instance, err := vm.New(program)
	if err != nil {
		t.Fatalf("create VM: %v", err)
	}

	t.Cleanup(func() { _ = instance.Close() })
	env, err := vm.NewEnvironment([]vm.EnvironmentOption{
		vm.WithNamespace(spec.Stdlib()),
		vm.WithParams(params),
	})
	if err != nil {
		t.Fatalf("create environment: %v", err)
	}

	result, err := instance.Run(context.Background(), env)
	if err != nil {
		t.Fatalf("expected ON ERROR recovery, got %v", err)
	}

	defer func() { _ = result.Close() }()
	if got := result.Root(); got != runtime.NewString("recovered") {
		t.Fatalf("expected ON ERROR result, got %T %#v", got, got)
	}
}

func runStreamGroupPanic(
	t *testing.T,
	query string,
	params map[string]runtime.Value,
	policy vm.PanicPolicy,
) (recovered any, runErr error) {
	t.Helper()
	program, err := spec.Compile(query)
	if err != nil {
		t.Fatalf("compile grouped wait: %v", err)
	}

	instance, err := vm.NewWith(program, vm.WithPanicPolicy(policy))
	if err != nil {
		t.Fatalf("create VM: %v", err)
	}

	t.Cleanup(func() { _ = instance.Close() })
	env, err := vm.NewEnvironment([]vm.EnvironmentOption{
		vm.WithNamespace(spec.Stdlib()),
		vm.WithParams(params),
	})
	if err != nil {
		t.Fatalf("create environment: %v", err)
	}

	defer func() {
		recovered = recover()
	}()
	result, runErr := instance.Run(context.Background(), env)
	if result != nil {
		_ = result.Close()
	}

	return nil, runErr
}
