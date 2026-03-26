package vm_test

import (
	"context"
	"errors"
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
