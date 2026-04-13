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
			Message: "Division by zero",
			Contains: []string{
				"DivideByZero: Division by zero",
				":2:8",
				"attempt to divide by zero",
				"Hint: Ensure the denominator is non-zero before division",
				"Note: Add a conditional check before dividing",
			},
			NotContains: []string{"~"},
		}),
		spec.NewSpec(
			"LET obj = {}\nRETURN obj.foo.bar",
			"obj.fql",
		).Expect().ExecError(ShouldBeRuntimeError, &ExpectedRuntimeError{
			Message: "Cannot read property \"bar\" of None",
			Contains: []string{
				"TypeError: Cannot read property \"bar\" of None",
				"property access on None",
				"Hint: Use optional chaining (?.) or check for None before accessing a member",
			},
			NotContains: []string{"Caused by:"},
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

			if got, want := runtimeErr.Message, "Missing parameter"; got != want {
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

				if got, want := span.Label, "missing parameter"; got != want {
					t.Fatalf("unexpected main span label: got %q, want %q", got, want)
				}
			}

			if !mainSpanFound {
				t.Fatal("expected a main error span")
			}

			formatted := runtimeErr.Format()
			for _, needle := range []string{
				"UnresolvedSymbol: Missing parameter",
				" --> missing_param.fql:1:8",
				"RETURN @foo",
				"^^^^ missing parameter",
				"Hint: Provide all required parameters",
				"Caused by: missed parameter: @foo",
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
				"^^^^ missing parameter",
				"called from inner (#1)",
				"called from middle (#2)",
				"called from outer (#3)",
				"Note: VM stack: outer -> middle -> inner",
				"Caused by: missed parameter: @foo",
			} {
				if !strings.Contains(formatted, needle) {
					t.Fatalf("expected formatted runtime error to contain %q, got:\n%s", needle, formatted)
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

			if got, want := runtimeErr.Message, "Invalid argument 2 type"; got != want {
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

				if got, want := span.Label, "argument 2 type mismatch"; got != want {
					t.Fatalf("unexpected main span label: got %q, want %q", got, want)
				}
			}

			if !mainSpanFound {
				t.Fatal("expected a main error span")
			}

			if runtimeErr.Cause == nil {
				t.Fatal("expected nested runtime error cause")
			}

			if got, want := runtimeErr.Cause.Error(), "expected String or Int or Object, but got Array"; got != want {
				t.Fatalf("unexpected runtime error cause: got %q, want %q", got, want)
			}
		})
	}
}
