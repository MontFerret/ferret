package vm_test

import (
	"context"
	"errors"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/compiler"
	"github.com/MontFerret/ferret/v2/pkg/file"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

func forEachCompiledVM(t *testing.T, source string, fn func(t *testing.T, instance *vm.VM)) {
	t.Helper()

	forEachCompiledVMWithOptions(t, source, nil, fn)
}

func forEachCompiledVMWithOptions(t *testing.T, source string, vmOpts []vm.Option, fn func(t *testing.T, instance *vm.VM)) {
	t.Helper()

	cases := []struct {
		name  string
		level compiler.OptimizationLevel
	}{
		{name: "Opt O0", level: compiler.O0},
		{name: "Opt O1", level: compiler.O1},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			program, err := compiler.New(compiler.WithOptimizationLevel(tc.level)).
				Compile(file.NewSource(tc.name, source))
			if err != nil {
				t.Fatalf("compile failed: %v", err)
			}

			instance, err := vm.NewWith(program, vmOpts...)
			if err != nil {
				t.Fatalf("vm init failed: %v", err)
			}

			fn(t, instance)
		})
	}
}

func mustNewEnvironment(t *testing.T, opts ...vm.EnvironmentOption) *vm.Environment {
	t.Helper()

	env, err := vm.NewEnvironment(opts)
	if err != nil {
		t.Fatalf("environment build failed: %v", err)
	}

	return env
}

func assertRuntimeErrorMessage(t *testing.T, err error, want string) *vm.RuntimeError {
	t.Helper()

	if err == nil {
		t.Fatal("expected runtime error")
	}

	var rtErr *vm.RuntimeError
	if !errors.As(err, &rtErr) {
		t.Fatalf("expected runtime error, got %T", err)
	}

	if rtErr.Message != want {
		t.Fatalf("unexpected runtime error message: got %q, want %q", rtErr.Message, want)
	}

	return rtErr
}

func assertRuntimeValueEquals(t *testing.T, got, want runtime.Value) {
	t.Helper()

	if got != want {
		t.Fatalf("unexpected result: got %v, want %v", got, want)
	}
}

func assertRuntimeIntArrayValue(t *testing.T, got runtime.Value, want ...runtime.Int) {
	t.Helper()

	arr, ok := got.(*runtime.Array)
	if !ok {
		t.Fatalf("expected runtime.Array, got %T", got)
	}

	expected := make([]runtime.Value, len(want))
	for i, value := range want {
		expected[i] = value
	}

	if arr.Compare(runtime.NewArrayWith(expected...)) != 0 {
		t.Fatalf("unexpected array value: got %v, want %v", got, runtime.NewArrayWith(expected...))
	}
}

func runVM(t *testing.T, instance *vm.VM, env *vm.Environment) (runtime.Value, error) {
	t.Helper()

	result, err := instance.Run(context.Background(), env)
	if err != nil {
		return runtime.None, err
	}

	root := result.Root()
	if closeErr := result.Close(); closeErr != nil {
		t.Fatalf("expected result close to succeed, got %v", closeErr)
	}

	return root, nil
}
