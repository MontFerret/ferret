package vm_test

import (
	"context"
	"strings"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

func TestRunMissingParamPrecedesWarmupHostResolution(t *testing.T) {
	forEachCompiledVM(t, "RETURN MISSING_FN(@foo)", func(t *testing.T, instance *vm.VM) {
		_, err := runVM(t, instance, vm.NewDefaultEnvironment())
		assertRuntimeErrorMessage(t, err, "Missing parameter")
	})
}

func TestStrictWarmupFailsProtectedMissingHostCallForDefaultAndBuiltEnvironment(t *testing.T) {
	forEachCompiledVM(t, "RETURN MISSING_FN()?", func(t *testing.T, instance *vm.VM) {
		builtEnv := mustNewEnvironment(t)

		cases := []struct {
			env  *vm.Environment
			name string
		}{
			{name: "default", env: vm.NewDefaultEnvironment()},
			{name: "built", env: builtEnv},
		}

		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				_, err := runVM(t, instance, tc.env)
				assertRuntimeErrorMessage(t, err, "Unresolved function")
			})
		}
	})
}

func TestStrictWarmupFailsOnDeadCodeUnresolvedHostCall(t *testing.T) {
	forEachCompiledVM(t, "RETURN false ? MISSING_FN() : 1", func(t *testing.T, instance *vm.VM) {
		envWithDummy := mustNewEnvironment(t, vm.WithFunction("DUMMY", func(context.Context, ...runtime.Value) (runtime.Value, error) {
			return runtime.None, nil
		}))

		cases := []struct {
			env  *vm.Environment
			name string
		}{
			{name: "default", env: vm.NewDefaultEnvironment()},
			{name: "dummy", env: envWithDummy},
		}

		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				_, err := runVM(t, instance, tc.env)
				assertRuntimeErrorMessage(t, err, "Unresolved function")
			})
		}
	})
}

func TestStrictWarmupAggregatesMissingHostFunctions(t *testing.T) {
	forEachCompiledVM(t, `
LET a = MISSING_A()
LET b = MISSING_B()
RETURN a + b
`, func(t *testing.T, instance *vm.VM) {
		_, err := runVM(t, instance, vm.NewDefaultEnvironment())
		if err == nil {
			t.Fatal("expected warmup error")
		}

		formatted := diagnostics.Format(err)
		if got, want := strings.Count(formatted, "Unresolved function"), 2; got != want {
			t.Fatalf("unexpected unresolved function count: got %d, want %d\n%s", got, want, formatted)
		}
	})
}

func TestStrictWarmupReportsRepeatedUnresolvedHostFunctionPerCallsite(t *testing.T) {
	forEachCompiledVM(t, `
LET a = MISSING()
LET b = MISSING()
RETURN a + b
`, func(t *testing.T, instance *vm.VM) {
		_, err := runVM(t, instance, vm.NewDefaultEnvironment())
		if err == nil {
			t.Fatal("expected warmup error")
		}

		formatted := diagnostics.Format(err)
		if got, want := strings.Count(formatted, "Unresolved function"), 2; got != want {
			t.Fatalf("unexpected unresolved function count: got %d, want %d\n%s", got, want, formatted)
		}
	})
}

func TestStrictWarmupFailureIsRepeatableUntilEnvironmentFixed(t *testing.T) {
	forEachCompiledVM(t, "RETURN F()", func(t *testing.T, instance *vm.VM) {
		_, err := runVM(t, instance, vm.NewDefaultEnvironment())
		assertRuntimeErrorMessage(t, err, "Unresolved function")

		_, err = runVM(t, instance, vm.NewDefaultEnvironment())
		assertRuntimeErrorMessage(t, err, "Unresolved function")

		validEnv := mustNewEnvironment(t, vm.WithFunction("F", func(context.Context, ...runtime.Value) (runtime.Value, error) {
			return runtime.NewInt(7), nil
		}))

		out, err := runVM(t, instance, validEnv)
		if err != nil {
			t.Fatalf("expected successful run after env fix, got %v", err)
		}

		assertRuntimeValueEquals(t, out, runtime.NewInt(7))
	})
}

func TestResetDrainsLeakedFramesBetweenFailedRuns(t *testing.T) {
	forEachCompiledVM(t, `
FUNC inner() (
	RETURN 1 / 0
)

FUNC outer() (
	RETURN inner()
)

RETURN outer()
`, func(t *testing.T, instance *vm.VM) {
		runAndCount := func(label string) int {
			t.Helper()

			_, err := runVM(t, instance, vm.NewDefaultEnvironment())
			rtErr := assertRuntimeErrorMessage(t, err, "Division by zero")

			return strings.Count(rtErr.Format(), "called from")
		}

		first := runAndCount("first")
		second := runAndCount("second")
		if first != second {
			t.Fatalf("expected stable stack depth across repeated failed runs: first=%d second=%d", first, second)
		}
	})
}

func TestHostNilResultIsNormalizedToNone(t *testing.T) {
	forEachCompiledVM(t, "RETURN NIL_FN()", func(t *testing.T, instance *vm.VM) {
		env := mustNewEnvironment(t, vm.WithFunction("NIL_FN", func(context.Context, ...runtime.Value) (runtime.Value, error) {
			return nil, nil
		}))

		out, err := runVM(t, instance, env)
		if err != nil {
			t.Fatalf("unexpected run error: %v", err)
		}

		if out != runtime.None {
			t.Fatalf("expected runtime.None, got %v", out)
		}
	})
}

func TestModuloTypeErrorNotMisclassifiedAsModuloByZero(t *testing.T) {
	forEachCompiledVM(t, `RETURN 5 % "x"`, func(t *testing.T, instance *vm.VM) {
		_, err := runVM(t, instance, vm.NewDefaultEnvironment())
		if err == nil {
			t.Fatal("expected runtime error")
		}

		rtErr := assertRuntimeErrorMessage(t, err, "Invalid type")
		if rtErr.Kind != diagnostics.TypeError {
			t.Fatalf("unexpected error kind: got %s, want %s", rtErr.Kind, diagnostics.TypeError)
		}

		if strings.Contains(strings.ToLower(rtErr.Format()), "modulo by zero") {
			t.Fatalf("expected non-modulo classification, got:\n%s", rtErr.Format())
		}
	})
}

func TestRuntimeErrorIncludesUDFCallStackContext(t *testing.T) {
	forEachCompiledVM(t, `
FUNC inner() (
	RETURN @x.foo
)
FUNC middle() (
	LET value = inner()
	RETURN value
)
FUNC outer() (
	LET value = middle()
	RETURN value
)
RETURN outer()
`, func(t *testing.T, instance *vm.VM) {
		env := vm.NewDefaultEnvironment()
		env.Params["x"] = runtime.None

		_, err := runVM(t, instance, env)
		rtErr := assertRuntimeErrorMessage(t, err, "Cannot read property \"foo\" of None")

		formatted := rtErr.Format()
		if !strings.Contains(formatted, "called from inner (#1)") {
			t.Fatalf("expected VM call stack context, got:\n%s", formatted)
		}

		if !strings.Contains(formatted, "VM stack: outer -> middle -> inner") {
			t.Fatalf("expected additive VM stack note, got:\n%s", formatted)
		}
	})
}

func TestRuntimeErrorSingleUdfStackFormattingUsesSourceSpelling(t *testing.T) {
	forEachCompiledVM(t, `
FUNC boo() (
	LET a = 1
	LET b = 0
	RETURN a / b
)
RETURN boo()
`, func(t *testing.T, instance *vm.VM) {
		_, err := runVM(t, instance, vm.NewDefaultEnvironment())
		rtErr := assertRuntimeErrorMessage(t, err, "Division by zero")

		formatted := rtErr.Format()
		if !strings.Contains(formatted, "called from boo (#1)") {
			t.Fatalf("expected call-site label with source-spelling udf name, got:\n%s", formatted)
		}

		if !strings.Contains(formatted, "VM stack: boo") {
			t.Fatalf("expected source-spelling VM stack note, got:\n%s", formatted)
		}
	})
}
