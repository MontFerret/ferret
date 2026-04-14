package vm_test

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/compiler"
	pkgdiagnostics "github.com/MontFerret/ferret/v2/pkg/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/source"
	"github.com/MontFerret/ferret/v2/pkg/vm"
	"github.com/MontFerret/ferret/v2/test/spec"
	. "github.com/MontFerret/ferret/v2/test/spec/exec"
)

func TestRuntimeErrorFormatting(t *testing.T) {
	RunSpecs(t, []spec.Spec{
		spec.NewSpec(
			"LET numerator = 10\nRETURN numerator / 0",
			"script.fql",
		).Expect().ExecError(ShouldBeRuntimeError, &ExpectedRuntimeError{
			Message: "division by zero",
			Contains: []string{
				"DivideByZero: division by zero",
				":2:8",
				"attempt to divide by zero",
				"Hint: Ensure the denominator is non-zero before division",
			},
			NotContains: []string{"~"},
		}),
		spec.NewSpec(
			"LET obj = {}\nRETURN obj.foo.bar",
			"obj.fql",
		).Expect().ExecError(ShouldBeRuntimeError, &ExpectedRuntimeError{
			Message: "invalid type",
			Contains: []string{
				"TypeError: invalid type",
				"property \"bar\" cannot be read from this value",
				"Note: cannot read property \"bar\" of None",
				"Hint: Use optional chaining (?.) or check for None before accessing a member",
				"Caused by: invalid type",
			},
		}),
		spec.NewSpec(
			`
FUNC Inner() => FAIL()
FUNC Outer() (
  LET result = Inner()
  RETURN result
)
RETURN Outer()
`,
			"nested_udf_stack.fql",
		).Expect().ExecError(ShouldBeRuntimeError, &ExpectedRuntimeError{
			Contains: []string{
				"called from Inner (#1)",
				"called from Outer (#2)",
				"VM stack: Outer -> Inner",
			},
		}).Env(
			vm.WithFunction("FAIL", func(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
				return runtime.None, errors.New("boom")
			}),
		),
	})
}

func TestRuntimeErrorFormatsMissingParamWithParamSpan(t *testing.T) {
	const query = `RETURN @foo`

	for _, level := range []compiler.OptimizationLevel{compiler.O0, compiler.O1} {
		t.Run(fmt.Sprintf("O%d", level), func(t *testing.T) {
			program, err := compiler.New(compiler.WithOptimizationLevel(level)).Compile(source.New("missing_param.fql", query))
			if err != nil {
				t.Fatalf("compile failed: %v", err)
			}

			instance, err := vm.New(program)
			if err != nil {
				t.Fatalf("vm init failed: %v", err)
			}
			defer func() {
				if closeErr := instance.Close(); closeErr != nil {
					t.Fatalf("vm close failed: %v", closeErr)
				}
			}()

			_, err = instance.Run(context.Background(), vm.NewDefaultEnvironment())
			if err == nil {
				t.Fatal("expected runtime error")
			}

			var runtimeErr *vm.RuntimeError
			if !errors.As(err, &runtimeErr) {
				t.Fatalf("expected runtime error, got %T", err)
			}

			if got, want := runtimeErr.Message, "missing parameter"; got != want {
				t.Fatalf("unexpected runtime error message: got %q, want %q", got, want)
			}

			mainSpanFound := false
			for _, span := range runtimeErr.Spans {
				if !span.Main {
					continue
				}

				mainSpanFound = true

				if got, want := query[span.Span.Start:span.Span.End], "@foo"; got != want {
					t.Fatalf("unexpected main span fragment: got %q, want %q", got, want)
				}

				if got, want := span.Label, "parameter '@foo' was not provided"; got != want {
					t.Fatalf("unexpected main span label: got %q, want %q", got, want)
				}
			}

			if !mainSpanFound {
				t.Fatal("expected a main error span")
			}

			formatted := runtimeErr.Format()
			for _, needle := range []string{
				"UnresolvedSymbol: missing parameter",
				" --> missing_param.fql:1:8",
				"RETURN @foo",
				"^^^^ parameter '@foo' was not provided",
				"Note: this query requires parameter '@foo'",
				"Hint: Provide a value for @foo before executing this query",
				"Caused by: missed parameter",
			} {
				if !strings.Contains(formatted, needle) {
					t.Fatalf("expected formatted runtime error to contain %q, got:\n%s", needle, formatted)
				}
			}
		})
	}
}

func TestRuntimeErrorFormatsNestedMissingParamAtInnerUsage(t *testing.T) {
	const query = `FUNC inner() => @foo
FUNC middle() (
  LET value = inner()
  RETURN value
)
FUNC outer() (
  LET value = middle()
  RETURN value
)
RETURN outer()
`

	for _, level := range []compiler.OptimizationLevel{compiler.O0, compiler.O1} {
		t.Run(fmt.Sprintf("O%d", level), func(t *testing.T) {
			program, err := compiler.New(compiler.WithOptimizationLevel(level)).Compile(source.New("missing_param_udf.fql", query))
			if err != nil {
				t.Fatalf("compile failed: %v", err)
			}

			instance, err := vm.New(program)
			if err != nil {
				t.Fatalf("vm init failed: %v", err)
			}
			defer func() {
				if closeErr := instance.Close(); closeErr != nil {
					t.Fatalf("vm close failed: %v", closeErr)
				}
			}()

			_, err = instance.Run(context.Background(), vm.NewDefaultEnvironment())
			if err == nil {
				t.Fatal("expected runtime error")
			}

			var runtimeErr *vm.RuntimeError
			if !errors.As(err, &runtimeErr) {
				t.Fatalf("expected runtime error, got %T", err)
			}

			mainSpanFound := false
			for _, span := range runtimeErr.Spans {
				if !span.Main {
					continue
				}

				mainSpanFound = true

				if got, want := query[span.Span.Start:span.Span.End], "@foo"; got != want {
					t.Fatalf("unexpected main span fragment: got %q, want %q", got, want)
				}
			}

			if !mainSpanFound {
				t.Fatal("expected a main error span")
			}

			formatted := runtimeErr.Format()
			for _, needle := range []string{
				" --> missing_param_udf.fql:1:17",
				"FUNC inner() => @foo",
				"^^^^ parameter '@foo' was not provided",
				"Note: this query requires parameter '@foo'",
				"Caused by: missed parameter",
			} {
				if !strings.Contains(formatted, needle) {
					t.Fatalf("expected formatted runtime error to contain %q, got:\n%s", needle, formatted)
				}
			}

			for _, needle := range []string{
				"called from",
				"VM stack:",
				"RETURN outer()",
				"RETURN value",
			} {
				if strings.Contains(formatted, needle) {
					t.Fatalf("expected formatted runtime error to not contain %q, got:\n%s", needle, formatted)
				}
			}
		})
	}
}

func TestRuntimeErrorFormatsAggregatedMissingParams(t *testing.T) {
	const query = `RETURN @foo + @bar`

	for _, level := range []compiler.OptimizationLevel{compiler.O0, compiler.O1} {
		t.Run(fmt.Sprintf("O%d", level), func(t *testing.T) {
			program, err := compiler.New(compiler.WithOptimizationLevel(level)).Compile(source.New("missing_params.fql", query))
			if err != nil {
				t.Fatalf("compile failed: %v", err)
			}

			instance, err := vm.New(program)
			if err != nil {
				t.Fatalf("vm init failed: %v", err)
			}
			defer func() {
				if closeErr := instance.Close(); closeErr != nil {
					t.Fatalf("vm close failed: %v", closeErr)
				}
			}()

			_, err = instance.Run(context.Background(), vm.NewDefaultEnvironment())
			if err == nil {
				t.Fatal("expected aggregated runtime error")
			}

			var runtimeErr *vm.RuntimeError
			if errors.As(err, &runtimeErr) {
				t.Fatalf("expected aggregated runtime error, got single runtime error: %v", runtimeErr)
			}

			formatted := pkgdiagnostics.Format(err)
			for _, needle := range []string{
				" --> missing_params.fql:1:8",
				" --> missing_params.fql:1:15",
				"parameter '@foo' was not provided",
				"parameter '@bar' was not provided",
			} {
				if !strings.Contains(formatted, needle) {
					t.Fatalf("expected formatted runtime error set to contain %q, got:\n%s", needle, formatted)
				}
			}

			if got, want := strings.Count(formatted, "UnresolvedSymbol: missing parameter"), 2; got != want {
				t.Fatalf("unexpected missing parameter diagnostic count: got %d, want %d\n%s", got, want, formatted)
			}
		})
	}
}

func TestRuntimeErrorFormatsAggregatedRepeatedMissingParamCallsites(t *testing.T) {
	const query = `RETURN @foo + @foo`

	for _, level := range []compiler.OptimizationLevel{compiler.O0, compiler.O1} {
		t.Run(fmt.Sprintf("O%d", level), func(t *testing.T) {
			program, err := compiler.New(compiler.WithOptimizationLevel(level)).Compile(source.New("missing_param_repeated.fql", query))
			if err != nil {
				t.Fatalf("compile failed: %v", err)
			}

			instance, err := vm.New(program)
			if err != nil {
				t.Fatalf("vm init failed: %v", err)
			}
			defer func() {
				if closeErr := instance.Close(); closeErr != nil {
					t.Fatalf("vm close failed: %v", closeErr)
				}
			}()

			_, err = instance.Run(context.Background(), vm.NewDefaultEnvironment())
			if err == nil {
				t.Fatal("expected aggregated runtime error")
			}

			formatted := pkgdiagnostics.Format(err)
			for _, needle := range []string{
				" --> missing_param_repeated.fql:1:8",
				" --> missing_param_repeated.fql:1:15",
			} {
				if !strings.Contains(formatted, needle) {
					t.Fatalf("expected formatted runtime error set to contain %q, got:\n%s", needle, formatted)
				}
			}

			if got, want := strings.Count(formatted, "UnresolvedSymbol: missing parameter"), 2; got != want {
				t.Fatalf("unexpected missing parameter diagnostic count: got %d, want %d\n%s", got, want, formatted)
			}

			if got, want := strings.Count(formatted, "parameter '@foo' was not provided"), 2; got != want {
				t.Fatalf("unexpected repeated missing parameter label count: got %d, want %d\n%s", got, want, formatted)
			}
		})
	}
}

func TestRuntimeErrorFormatsAggregatedUdfMissingParamCallsites(t *testing.T) {
	const query = `FUNC read() => @foo
LET left = read()
LET right = read()
RETURN left + right
`

	for _, level := range []compiler.OptimizationLevel{compiler.O0, compiler.O1} {
		t.Run(fmt.Sprintf("O%d", level), func(t *testing.T) {
			program, err := compiler.New(compiler.WithOptimizationLevel(level)).Compile(source.New("missing_param_udf_callsites.fql", query))
			if err != nil {
				t.Fatalf("compile failed: %v", err)
			}

			instance, err := vm.New(program)
			if err != nil {
				t.Fatalf("vm init failed: %v", err)
			}
			defer func() {
				if closeErr := instance.Close(); closeErr != nil {
					t.Fatalf("vm close failed: %v", closeErr)
				}
			}()

			_, err = instance.Run(context.Background(), vm.NewDefaultEnvironment())
			if err == nil {
				t.Fatal("expected aggregated runtime error")
			}

			formatted := pkgdiagnostics.Format(err)
			for _, needle := range []string{
				" --> missing_param_udf_callsites.fql:1:16",
			} {
				if !strings.Contains(formatted, needle) {
					t.Fatalf("expected formatted runtime error set to contain %q, got:\n%s", needle, formatted)
				}
			}

			if got, want := strings.Count(formatted, "UnresolvedSymbol: missing parameter"), 1; got != want {
				t.Fatalf("unexpected missing parameter diagnostic count: got %d, want %d\n%s", got, want, formatted)
			}

			if got, want := strings.Count(formatted, "parameter '@foo' was not provided"), 1; got != want {
				t.Fatalf("unexpected missing parameter label count: got %d, want %d\n%s", got, want, formatted)
			}

			for _, needle := range []string{
				"called from",
				"VM stack:",
			} {
				if strings.Contains(formatted, needle) {
					t.Fatalf("expected formatted runtime error set to not contain %q, got:\n%s", needle, formatted)
				}
			}
		})
	}
}

func TestRuntimeErrorFormatsAggregatedTopLevelAndUdfMissingParams(t *testing.T) {
	const query = `LET val = @foo
LET val2 = @bar

FUNC TEST() (
  RETURN @baz
)

RETURN [val, val2, TEST()]
`

	for _, level := range []compiler.OptimizationLevel{compiler.O0, compiler.O1} {
		t.Run(fmt.Sprintf("O%d", level), func(t *testing.T) {
			program, err := compiler.New(compiler.WithOptimizationLevel(level)).Compile(source.New("missing_param_mixed_sites.fql", query))
			if err != nil {
				t.Fatalf("compile failed: %v", err)
			}

			instance, err := vm.New(program)
			if err != nil {
				t.Fatalf("vm init failed: %v", err)
			}
			defer func() {
				if closeErr := instance.Close(); closeErr != nil {
					t.Fatalf("vm close failed: %v", closeErr)
				}
			}()

			_, err = instance.Run(context.Background(), vm.NewDefaultEnvironment())
			if err == nil {
				t.Fatal("expected aggregated runtime error")
			}

			formatted := pkgdiagnostics.Format(err)
			for _, needle := range []string{
				"LET val = @foo",
				"LET val2 = @bar",
				"RETURN @baz",
				"parameter '@foo' was not provided",
				"parameter '@bar' was not provided",
				"parameter '@baz' was not provided",
			} {
				if !strings.Contains(formatted, needle) {
					t.Fatalf("expected formatted runtime error set to contain %q, got:\n%s", needle, formatted)
				}
			}

			if got, want := strings.Count(formatted, "UnresolvedSymbol: missing parameter"), 3; got != want {
				t.Fatalf("unexpected missing parameter diagnostic count: got %d, want %d\n%s", got, want, formatted)
			}

			for _, needle := range []string{
				"called from",
				"VM stack:",
			} {
				if strings.Contains(formatted, needle) {
					t.Fatalf("expected formatted runtime error set to not contain %q, got:\n%s", needle, formatted)
				}
			}
		})
	}
}

func TestRuntimeErrorFormatsArgumentTypeFailuresWithArgumentSpan(t *testing.T) {
	const query = `RETURN BROKEN("ok", [1, 2])`

	for _, level := range []compiler.OptimizationLevel{compiler.O0, compiler.O1} {
		t.Run(fmt.Sprintf("O%d", level), func(t *testing.T) {
			program, err := compiler.New(compiler.WithOptimizationLevel(level)).Compile(source.New("arg_type.fql", query))
			if err != nil {
				t.Fatalf("compile failed: %v", err)
			}

			instance, err := vm.New(program)
			if err != nil {
				t.Fatalf("vm init failed: %v", err)
			}
			defer func() {
				if closeErr := instance.Close(); closeErr != nil {
					t.Fatalf("vm close failed: %v", closeErr)
				}
			}()

			env, err := vm.NewEnvironment([]vm.EnvironmentOption{
				vm.WithFunction("BROKEN", func(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
					return runtime.None, runtime.ValidateArgTypeAt(args, 1, runtime.TypeString, runtime.TypeInt, runtime.TypeObject)
				}),
			})
			if err != nil {
				t.Fatalf("environment init failed: %v", err)
			}

			_, err = instance.Run(context.Background(), env)
			if err == nil {
				t.Fatal("expected runtime error")
			}

			var runtimeErr *vm.RuntimeError
			if !errors.As(err, &runtimeErr) {
				t.Fatalf("expected runtime error, got %T", err)
			}

			if got, want := runtimeErr.Message, "invalid argument type"; got != want {
				t.Fatalf("unexpected runtime error message: got %q, want %q", got, want)
			}

			if got, want := runtimeErr.Kind, pkgdiagnostics.TypeError; got != want {
				t.Fatalf("unexpected runtime error kind: got %s, want %s", got, want)
			}

			mainSpanFound := false
			for _, span := range runtimeErr.Spans {
				if !span.Main {
					continue
				}

				mainSpanFound = true

				if got, want := query[span.Span.Start:span.Span.End], "[1, 2]"; got != want {
					t.Fatalf("unexpected main span fragment: got %q, want %q", got, want)
				}

				if got, want := span.Label, "argument 2 has incompatible type"; got != want {
					t.Fatalf("unexpected main span label: got %q, want %q", got, want)
				}
			}

			if !mainSpanFound {
				t.Fatal("expected a main error span")
			}

			if runtimeErr.Cause == nil {
				t.Fatal("expected nested runtime error cause")
			}

			if got, want := runtimeErr.Note, "argument 2 expects String or Int or Object, but got Array"; got != want {
				t.Fatalf("unexpected runtime error note: got %q, want %q", got, want)
			}

			if got, want := runtimeErr.Hint, "Convert argument 2 to String or Int or Object before this call"; got != want {
				t.Fatalf("unexpected runtime error hint: got %q, want %q", got, want)
			}

			if !strings.Contains(runtimeErr.Format(), "Hint: Convert argument 2 to String or Int or Object before this call") {
				t.Fatalf("expected formatted runtime error to include synthesized hint, got:\n%s", runtimeErr.Format())
			}

			if got, want := runtimeErr.Cause.Error(), "invalid type"; got != want {
				t.Fatalf("unexpected runtime error cause: got %q, want %q", got, want)
			}
		})
	}
}
