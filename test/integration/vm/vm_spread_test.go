package vm_test

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/compiler"
	pkgdiagnostics "github.com/MontFerret/ferret/v2/pkg/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/source"
	"github.com/MontFerret/ferret/v2/pkg/vm"
	"github.com/MontFerret/ferret/v2/test/spec"
	. "github.com/MontFerret/ferret/v2/test/spec/exec"
)

func TestLiteralSpread(t *testing.T) {
	RunSpecs(t, []spec.Spec{
		Array(`RETURN [0, ...[1, 2], 3, ...[], ...NONE]`, []any{0, 1, 2, 3}, "array spread preserves entry order"),
		Object(`RETURN { a: 0, ...{ a: 1, b: 2 }, a: 3, ...{ c: 4 }, ...NONE }`, map[string]any{
			"a": 3,
			"b": 2,
			"c": 4,
		}, "object spread preserves later-property precedence"),
		Array(`
LET source = [1]
LET copied = [...source]
copied[0] = 9
RETURN [source, copied]
`, []any{[]any{1}, []any{9}}, "array spread does not mutate its source container"),
		Array(`
LET nested = { value: 1 }
LET source = [nested]
LET copied = [...source]
copied[0].value = 2
RETURN [source[0].value, copied[0].value]
`, []any{2, 2}, "array spread is shallow"),
		Array(`
LET source = { value: 1 }
LET copied = { ...source }
copied.value = 9
RETURN [source.value, copied.value]
`, []any{1, 9}, "object spread does not mutate its source container"),
		Array(`
LET nested = { value: 1 }
LET source = { nested }
LET copied = { ...source }
copied.nested.value = 2
RETURN [source.nested.value, copied.nested.value]
`, []any{2, 2}, "object spread is shallow"),
		Array(`RETURN [...HOST_ARRAY()]`, []any{1, 2}, "host-returned Array can be spread"),
		Object(`RETURN { ...HOST_OBJECT() }`, map[string]any{"value": 1}, "host-returned Object can be spread"),
		Error(`RETURN [...(1..2)]`, "generic List values cannot be spread into arrays"),
		Error(`RETURN { ...HOST_MAP() }`, "Map values without ObjectLike cannot be spread into objects"),
		Array(`RETURN ([...FAIL()]) ON ERROR RETURN ["recovered"]`, []any{"recovered"}, "spread-source failures use normal recovery"),
	},
		vm.WithFunction("HOST_ARRAY", func(context.Context, ...runtime.Value) (runtime.Value, error) {
			return runtime.NewArrayWith(runtime.NewInt(1), runtime.NewInt(2)), nil
		}),
		vm.WithFunction("HOST_OBJECT", func(context.Context, ...runtime.Value) (runtime.Value, error) {
			return runtime.NewObjectWith(map[string]runtime.Value{"value": runtime.NewInt(1)}), nil
		}),
		vm.WithFunction("HOST_MAP", func(context.Context, ...runtime.Value) (runtime.Value, error) {
			return genericMap{Map: runtime.NewObjectWith(map[string]runtime.Value{"value": runtime.NewInt(1)})}, nil
		}),
		vm.WithFunction("FAIL", func(context.Context, ...runtime.Value) (runtime.Value, error) {
			return runtime.None, errors.New("spread source failed")
		}),
	)
}

func TestLiteralSpreadEvaluationOrder(t *testing.T) {
	query := `RETURN [
  TRACE(1, 10),
  ...TRACE(2, [20, 30]),
  TRACE(3, 40),
  { first: TRACE(4, 1), ...TRACE(5, { second: 2 }), third: TRACE(6, 3) }
]`

	for _, level := range []compiler.OptimizationLevel{compiler.O0, compiler.O1} {
		t.Run(fmt.Sprintf("O%d", level), func(t *testing.T) {
			var trace []int

			RunSpecsWith(t, "trace", compiler.New(compiler.WithOptimizationLevel(level)), []spec.Spec{
				Array(query, []any{10, 20, 30, 40, map[string]any{
					"first":  1,
					"second": 2,
					"third":  3,
				}}),
			}, vm.WithFunction("TRACE", func(_ context.Context, args ...runtime.Value) (runtime.Value, error) {
				trace = append(trace, int(args[0].(runtime.Int)))

				return args[1], nil
			}))

			if want := []int{1, 2, 3, 4, 5, 6}; !slices.Equal(trace, want) {
				t.Fatalf("unexpected evaluation trace: got %v, want %v", trace, want)
			}
		})
	}
}

func TestLiteralSpreadTypeErrors(t *testing.T) {
	tests := []struct {
		query   string
		message string
		span    string
		hint    string
	}{
		{
			query:   `RETURN [..."value"]`,
			message: "cannot spread String into Array",
			span:    `..."value"`,
			hint:    "Spread an Array value or none inside an array literal",
		},
		{
			query:   `RETURN {...1}`,
			message: "cannot spread Int into Object",
			span:    `...1`,
			hint:    "Spread an Object value or none inside an object literal",
		},
	}

	for _, level := range []compiler.OptimizationLevel{compiler.O0, compiler.O1} {
		for _, test := range tests {
			t.Run(fmt.Sprintf("O%d/%s", level, test.message), func(t *testing.T) {
				program, err := compiler.New(compiler.WithOptimizationLevel(level)).Compile(source.New("spread.fql", test.query))
				if err != nil {
					t.Fatalf("compile failed: %v", err)
				}

				instance, err := vm.New(program)
				if err != nil {
					t.Fatalf("vm init failed: %v", err)
				}
				defer instance.Close()

				env, err := vm.NewEnvironment(nil)
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

				if got := runtimeErr.Message; got != test.message {
					t.Fatalf("unexpected message: got %q, want %q", got, test.message)
				}

				if got := runtimeErr.Kind; got != pkgdiagnostics.TypeError {
					t.Fatalf("unexpected kind: got %s, want %s", got, pkgdiagnostics.TypeError)
				}

				if got := runtimeErr.Hint; got != test.hint {
					t.Fatalf("unexpected hint: got %q, want %q", got, test.hint)
				}

				if got := runtimeErr.Note; got != test.message {
					t.Fatalf("unexpected note: got %q, want %q", got, test.message)
				}

				if !errors.Is(runtimeErr, runtime.ErrInvalidType) {
					t.Fatalf("expected invalid-type cause, got %v", runtimeErr)
				}

				assertSpreadErrorSpan(t, runtimeErr, test.query, test.span)
			})
		}
	}
}

func assertSpreadErrorSpan(t *testing.T, runtimeErr *vm.RuntimeError, query, want string) {
	t.Helper()

	for _, span := range runtimeErr.Spans {
		if !span.Main {
			continue
		}

		if got := query[span.Span.Start:span.Span.End]; got != want {
			t.Fatalf("unexpected main span: got %q, want %q", got, want)
		}

		if !slices.Contains([]string{
			"expected Array or none in array literal",
			"expected Object or none in object literal",
		}, span.Label) {
			t.Fatalf("unexpected spread error label: %q", span.Label)
		}

		return
	}

	t.Fatal("expected a main spread error span")
}
