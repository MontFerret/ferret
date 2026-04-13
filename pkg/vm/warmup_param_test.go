package vm

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/compiler"
	pkgdiagnostics "github.com/MontFerret/ferret/v2/pkg/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/source"
	rtdiagnostics "github.com/MontFerret/ferret/v2/pkg/vm/internal/diagnostics"
)

func mustCompileProgram(t *testing.T, level compiler.OptimizationLevel, name, query string) *bytecode.Program {
	t.Helper()

	program, err := compiler.New(compiler.WithOptimizationLevel(level)).Compile(source.New(name, query))
	if err != nil {
		t.Fatalf("compile failed: %v", err)
	}

	return program
}

func TestWarmupMissingParamsAggregateTopLevelAndUdfSites(t *testing.T) {
	const query = `LET val = @foo
LET val2 = @bar

FUNC TEST() (
  RETURN @baz
)

RETURN [val, val2, TEST()]
`

	for _, level := range []compiler.OptimizationLevel{compiler.O0, compiler.O1} {
		t.Run(fmt.Sprintf("O%d", level), func(t *testing.T) {
			program := mustCompileProgram(t, level, "missing_param_mixed_sites.fql", query)
			instance := mustNewVM(t, program)

			if got, want := len(instance.plan.paramLoadDescriptors), 3; got != want {
				t.Fatalf("unexpected param descriptor count: got %d, want %d", got, want)
			}

			_, err := instance.Run(context.Background(), NewDefaultEnvironment())
			if err == nil {
				t.Fatal("expected aggregated runtime error")
			}

			var rtErrSet *rtdiagnostics.RuntimeErrorSet
			if !errors.As(err, &rtErrSet) {
				t.Fatalf("expected runtime error set, got %T", err)
			}

			if got, want := rtErrSet.Size(), 3; got != want {
				t.Fatalf("unexpected runtime error set size: got %d, want %d", got, want)
			}

			formatted := pkgdiagnostics.Format(err)
			for _, needle := range []string{
				"LET val = @foo",
				"LET val2 = @bar",
				"RETURN @baz",
			} {
				if !strings.Contains(formatted, needle) {
					t.Fatalf("expected formatted error to contain %q, got:\n%s", needle, formatted)
				}
			}

			for _, needle := range []string{
				"called from",
				"VM stack:",
			} {
				if strings.Contains(formatted, needle) {
					t.Fatalf("expected formatted error to not contain %q, got:\n%s", needle, formatted)
				}
			}
		})
	}
}

func TestWarmupMissingParamNestedUdfUsesOnlyUdfBodySnippet(t *testing.T) {
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
			program := mustCompileProgram(t, level, "missing_param_nested_udf.fql", query)
			instance := mustNewVM(t, program)

			_, err := instance.Run(context.Background(), NewDefaultEnvironment())
			if err == nil {
				t.Fatal("expected runtime error")
			}

			var rtErr *RuntimeError
			if !errors.As(err, &rtErr) {
				t.Fatalf("expected runtime error, got %T", err)
			}

			formatted := rtErr.Format()
			for _, needle := range []string{
				"FUNC inner() => @foo",
				"^^^^ missing parameter",
			} {
				if !strings.Contains(formatted, needle) {
					t.Fatalf("expected formatted error to contain %q, got:\n%s", needle, formatted)
				}
			}

			for _, needle := range []string{
				"called from",
				"VM stack:",
				"RETURN outer()",
				"RETURN value",
			} {
				if strings.Contains(formatted, needle) {
					t.Fatalf("expected formatted error to not contain %q, got:\n%s", needle, formatted)
				}
			}
		})
	}
}

func TestWarmupMissingParamRepeatedUdfCallsStillReportSingleLoadSite(t *testing.T) {
	const query = `FUNC read() => @foo
LET left = read()
LET right = read()
RETURN left + right
`

	for _, level := range []compiler.OptimizationLevel{compiler.O0, compiler.O1} {
		t.Run(fmt.Sprintf("O%d", level), func(t *testing.T) {
			program := mustCompileProgram(t, level, "missing_param_udf_callsites.fql", query)
			instance := mustNewVM(t, program)

			if got, want := len(instance.plan.paramLoadDescriptors), 1; got != want {
				t.Fatalf("unexpected param descriptor count: got %d, want %d", got, want)
			}

			_, err := instance.Run(context.Background(), NewDefaultEnvironment())
			if err == nil {
				t.Fatal("expected runtime error")
			}

			var rtErr *RuntimeError
			if !errors.As(err, &rtErr) {
				t.Fatalf("expected single runtime error, got %T", err)
			}

			formatted := rtErr.Format()
			if got, want := strings.Count(formatted, "Missing parameter"), 1; got != want {
				t.Fatalf("unexpected missing parameter count: got %d, want %d\n%s", got, want, formatted)
			}

			for _, needle := range []string{
				"called from",
				"VM stack:",
			} {
				if strings.Contains(formatted, needle) {
					t.Fatalf("expected formatted error to not contain %q, got:\n%s", needle, formatted)
				}
			}
		})
	}
}

func TestWarmupMissingParamProtectedUdfCallStillFailsWithoutTrace(t *testing.T) {
	const query = `FUNC risky() => @foo
RETURN risky()?
`

	for _, level := range []compiler.OptimizationLevel{compiler.O0, compiler.O1} {
		t.Run(fmt.Sprintf("O%d", level), func(t *testing.T) {
			program := mustCompileProgram(t, level, "missing_param_protected_udf.fql", query)
			instance := mustNewVM(t, program)

			_, err := instance.Run(context.Background(), NewDefaultEnvironment())
			if err == nil {
				t.Fatal("expected runtime error")
			}

			var rtErr *RuntimeError
			if !errors.As(err, &rtErr) {
				t.Fatalf("expected runtime error, got %T", err)
			}

			if got, want := rtErr.Message, "Missing parameter"; got != want {
				t.Fatalf("unexpected runtime error message: got %q, want %q", got, want)
			}

			formatted := rtErr.Format()
			for _, needle := range []string{
				"called from",
				"VM stack:",
			} {
				if strings.Contains(formatted, needle) {
					t.Fatalf("expected formatted error to not contain %q, got:\n%s", needle, formatted)
				}
			}
		})
	}
}
