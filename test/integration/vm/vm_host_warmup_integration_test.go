package vm_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

func TestWarmupClearsStaleHostCacheAcrossEnvironments(t *testing.T) {
	forEachCompiledVM(t, "RETURN F()", func(t *testing.T, instance *vm.VM) {
		envWithFn := mustNewEnvironment(t, vm.WithFunction("F", func(context.Context, ...runtime.Value) (runtime.Value, error) {
			return runtime.NewInt(7), nil
		}))

		out, err := runVM(t, instance, envWithFn)
		if err != nil {
			t.Fatalf("first run failed: %v", err)
		}
		assertRuntimeValueEquals(t, out, runtime.NewInt(7))

		_, err = runVM(t, instance, vm.NewDefaultEnvironment())
		assertRuntimeErrorMessage(t, err, "Unresolved function")
	})
}

func TestWarmupRebindsWhenEnvironmentFunctionNamesMatch(t *testing.T) {
	forEachCompiledVM(t, "RETURN F()", func(t *testing.T, instance *vm.VM) {
		envA := mustNewEnvironment(t, vm.WithFunction("F", func(context.Context, ...runtime.Value) (runtime.Value, error) {
			return runtime.NewInt(1), nil
		}))
		envB := mustNewEnvironment(t, vm.WithFunction("F", func(context.Context, ...runtime.Value) (runtime.Value, error) {
			return runtime.NewInt(2), nil
		}))

		out, err := runVM(t, instance, envA)
		if err != nil {
			t.Fatalf("first run failed: %v", err)
		}
		assertRuntimeValueEquals(t, out, runtime.NewInt(1))

		out, err = runVM(t, instance, envB)
		if err != nil {
			t.Fatalf("second run failed: %v", err)
		}
		assertRuntimeValueEquals(t, out, runtime.NewInt(2))
	})
}

func TestWarmupRebindThrashesAcrossSameNamedHostImplementations(t *testing.T) {
	forEachCompiledVM(t, "RETURN F()", func(t *testing.T, instance *vm.VM) {
		envA := mustNewEnvironment(t, vm.WithFunctionsRegistrar(func(fns runtime.FunctionDefs) {
			fns.A0().Add("F", func(context.Context) (runtime.Value, error) {
				return runtime.NewInt(1), nil
			})
		}))
		envB := mustNewEnvironment(t, vm.WithFunctionsRegistrar(func(fns runtime.FunctionDefs) {
			fns.A0().Add("F", func(context.Context) (runtime.Value, error) {
				return runtime.NewInt(2), nil
			})
		}))

		out, err := runVM(t, instance, envA)
		if err != nil {
			t.Fatalf("first run failed: %v", err)
		}
		assertRuntimeValueEquals(t, out, runtime.NewInt(1))

		out, err = runVM(t, instance, envB)
		if err != nil {
			t.Fatalf("second run failed: %v", err)
		}
		assertRuntimeValueEquals(t, out, runtime.NewInt(2))

		out, err = runVM(t, instance, envA)
		if err != nil {
			t.Fatalf("third run failed: %v", err)
		}
		assertRuntimeValueEquals(t, out, runtime.NewInt(1))
	})
}

func TestWarmupRebindSwitchesBetweenFixedArityAndVarargImplementations(t *testing.T) {
	forEachCompiledVM(t, "RETURN F(1, 2)", func(t *testing.T, instance *vm.VM) {
		envA2 := mustNewEnvironment(t, vm.WithFunctionsRegistrar(func(fns runtime.FunctionDefs) {
			fns.A2().Add("F", func(context.Context, runtime.Value, runtime.Value) (runtime.Value, error) {
				return runtime.NewInt(12), nil
			})
		}))
		envVar := mustNewEnvironment(t, vm.WithFunction("F", func(context.Context, ...runtime.Value) (runtime.Value, error) {
			return runtime.NewInt(102), nil
		}))
		envA2Again := mustNewEnvironment(t, vm.WithFunctionsRegistrar(func(fns runtime.FunctionDefs) {
			fns.A2().Add("F", func(context.Context, runtime.Value, runtime.Value) (runtime.Value, error) {
				return runtime.NewInt(22), nil
			})
		}))

		out, err := runVM(t, instance, envA2)
		if err != nil {
			t.Fatalf("fixed-arity run failed: %v", err)
		}
		assertRuntimeValueEquals(t, out, runtime.NewInt(12))

		out, err = runVM(t, instance, envVar)
		if err != nil {
			t.Fatalf("vararg run failed: %v", err)
		}
		assertRuntimeValueEquals(t, out, runtime.NewInt(102))

		out, err = runVM(t, instance, envA2Again)
		if err != nil {
			t.Fatalf("second fixed-arity run failed: %v", err)
		}
		assertRuntimeValueEquals(t, out, runtime.NewInt(22))
	})
}

func TestWarmupRepeatedMissingRunsAfterSuccessRecoverCleanly(t *testing.T) {
	forEachCompiledVM(t, "RETURN F()", func(t *testing.T, instance *vm.VM) {
		validEnv := mustNewEnvironment(t, vm.WithFunctionsRegistrar(func(fns runtime.FunctionDefs) {
			fns.A0().Add("F", func(context.Context) (runtime.Value, error) {
				return runtime.NewInt(7), nil
			})
		}))
		missingEnv := vm.NewDefaultEnvironment()
		recoveredEnv := mustNewEnvironment(t, vm.WithFunctionsRegistrar(func(fns runtime.FunctionDefs) {
			fns.A0().Add("F", func(context.Context) (runtime.Value, error) {
				return runtime.NewInt(9), nil
			})
		}))

		out, err := runVM(t, instance, validEnv)
		if err != nil {
			t.Fatalf("initial run failed: %v", err)
		}
		assertRuntimeValueEquals(t, out, runtime.NewInt(7))

		_, err = runVM(t, instance, missingEnv)
		assertRuntimeErrorMessage(t, err, "Unresolved function")

		_, err = runVM(t, instance, missingEnv)
		assertRuntimeErrorMessage(t, err, "Unresolved function")

		out, err = runVM(t, instance, recoveredEnv)
		if err != nil {
			t.Fatalf("recovery run failed: %v", err)
		}
		assertRuntimeValueEquals(t, out, runtime.NewInt(9))
	})
}

func TestWarmupMultiCallsiteRecoveryAfterPartialMissingRun(t *testing.T) {
	forEachCompiledVM(t, `
LET a = F()
LET b = G()
LET c = F()
RETURN [a, b, c]
`, func(t *testing.T, instance *vm.VM) {
		envAll := mustNewEnvironment(t, vm.WithFunctionsRegistrar(func(fns runtime.FunctionDefs) {
			fns.A0().Add("F", func(context.Context) (runtime.Value, error) {
				return runtime.NewInt(1), nil
			})
			fns.A0().Add("G", func(context.Context) (runtime.Value, error) {
				return runtime.NewInt(2), nil
			})
		}))
		missingEnv := mustNewEnvironment(t, vm.WithFunctionsRegistrar(func(fns runtime.FunctionDefs) {
			fns.A0().Add("F", func(context.Context) (runtime.Value, error) {
				return runtime.NewInt(20), nil
			})
		}))
		recoveredEnv := mustNewEnvironment(t, vm.WithFunctionsRegistrar(func(fns runtime.FunctionDefs) {
			fns.A0().Add("F", func(context.Context) (runtime.Value, error) {
				return runtime.NewInt(30), nil
			})
			fns.A0().Add("G", func(context.Context) (runtime.Value, error) {
				return runtime.NewInt(40), nil
			})
		}))

		out, err := runVM(t, instance, envAll)
		if err != nil {
			t.Fatalf("initial run failed: %v", err)
		}
		assertRuntimeIntArrayValue(t, out, runtime.NewInt(1), runtime.NewInt(2), runtime.NewInt(1))

		_, err = runVM(t, instance, missingEnv)
		assertRuntimeErrorMessage(t, err, "Unresolved function")

		out, err = runVM(t, instance, recoveredEnv)
		if err != nil {
			t.Fatalf("recovery run failed: %v", err)
		}
		assertRuntimeIntArrayValue(t, out, runtime.NewInt(30), runtime.NewInt(40), runtime.NewInt(30))
	})
}
