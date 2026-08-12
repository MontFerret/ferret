package vm_test

import (
	"context"
	"fmt"
	"reflect"
	"sync/atomic"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/compiler"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
	"github.com/MontFerret/ferret/v2/test/spec"
	. "github.com/MontFerret/ferret/v2/test/spec/exec"
)

func TestBodyCompletion(t *testing.T) {
	RunSpecs(t, []spec.Spec{
		Nil(`LET value = 1`, "effect-only script falls through"),
		Nil(`FUNC empty() {} RETURN empty()`, "empty block UDF falls through"),
		Nil(`FUNC value() { 42 } RETURN value()`, "final expression statement is discarded"),
		Nil(`
FOR value IN [1, 2, 3] {
  RETURN value * 2
}
`, "final standalone FOR result is discarded"),
		S(`
VAR total = 0
FOR value IN [1, 2] {
  total += value
  RETURN value
}
FOR value IN [3, 4] {
  total += value
  RETURN value
}
RETURN total
`, 10, "standalone FOR statements execute in source order"),
		Array(`
VAR total = 0
FUNC accumulate(items) {
  FUNC add(item) {
    total += item
  }
  FOR item IN items {
    add(item)
    RETURN item
  }
}
RETURN [accumulate([1, 2]), total]
`, []any{nil, 3}, "standalone UDF FOR preserves captures and reachable nested calls"),
		S(`FUNC value() => 42 RETURN value()`, 42, "arrow UDF remains value-producing"),
		S(`FUNC value() { RETURN 42 } RETURN value()`, 42, "explicit UDF return remains value-producing"),
		Nil(`RETURN NONE`, "explicit NONE remains value-producing"),
		Error(`
FOR value IN [1, 0] {
  RETURN 10 / value
}
`, "discarded FOR still propagates runtime errors"),
	})
}

func TestStandaloneForClosesDiscardedIterator(t *testing.T) {
	for _, level := range []compiler.OptimizationLevel{compiler.O0, compiler.O1} {
		t.Run(optimizationName(level), func(t *testing.T) {
			iterable := newBodyCompletionIterable()
			program, err := spec.Compile(`FOR value IN SOURCE() { RETURN value }`, level)
			if err != nil {
				t.Fatalf("compile query: %v", err)
			}

			out, err := spec.Run(program, vm.WithFunction("SOURCE", func(context.Context, ...runtime.Value) (runtime.Value, error) {
				return iterable, nil
			}))
			if err != nil {
				t.Fatalf("run query: %v", err)
			}

			if got, want := string(out), "null"; got != want {
				t.Fatalf("result = %s, want %s", got, want)
			}

			if got, want := iterable.closed(), int32(1); got != want {
				t.Fatalf("iterator close count = %d, want %d", got, want)
			}
		})
	}
}

func TestDiscardedForPreservesLoopSemantics(t *testing.T) {
	tests := []struct {
		name  string
		query string
		want  []int
	}{
		{
			name: "filter",
			query: `FOR value IN [1, 2, 3] {
  FILTER value > 1
  RETURN RECORD(value)
}`,
			want: []int{2, 3},
		},
		{
			name: "sort",
			query: `FOR value IN [3, 1, 2] {
  SORT value
  RETURN RECORD(value)
}`,
			want: []int{1, 2, 3},
		},
		{
			name: "limit",
			query: `FOR value IN [1, 2, 3] {
  LIMIT 2
  RETURN RECORD(value)
}`,
			want: []int{1, 2},
		},
		{
			name: "collect",
			query: `FOR value IN [1, 1, 2] {
  COLLECT key = value
  RETURN RECORD(key)
}`,
			want: []int{1, 2},
		},
		{
			name: "global aggregate",
			query: `FOR value IN [1, 2, 3] {
  COLLECT AGGREGATE total = COUNT(value)
  RETURN RECORD(total)
}`,
			want: []int{3},
		},
		{
			name: "pass-through nesting",
			query: `FOR outer IN [1, 2] {
  FOR inner IN [outer, outer + 1] {
    RETURN RECORD(inner)
  }
}`,
			want: []int{1, 2, 2, 3},
		},
	}

	for _, level := range []compiler.OptimizationLevel{compiler.O0, compiler.O1} {
		for _, test := range tests {
			t.Run(fmt.Sprintf("%s/%s", optimizationName(level), test.name), func(t *testing.T) {
				var calls []int
				program, err := spec.Compile(test.query, level)
				if err != nil {
					t.Fatalf("compile query: %v", err)
				}

				out, err := spec.Run(program, vm.WithFunction("RECORD", func(_ context.Context, args ...runtime.Value) (runtime.Value, error) {
					value := args[0].(runtime.Int)
					calls = append(calls, int(value))

					return value, nil
				}))
				if err != nil {
					t.Fatalf("run query: %v", err)
				}

				if got, want := string(out), "null"; got != want {
					t.Fatalf("result = %s, want %s", got, want)
				}

				if !reflect.DeepEqual(calls, test.want) {
					t.Fatalf("recorded values = %v, want %v", calls, test.want)
				}
			})
		}
	}
}

func TestDiscardedWrappedForRecovery(t *testing.T) {
	for _, level := range []compiler.OptimizationLevel{compiler.O0, compiler.O1} {
		t.Run(optimizationName(level), func(t *testing.T) {
			t.Run("suppress", func(t *testing.T) {
				calls := 0
				runDiscardedRecoveryQuery(t, level, `FUNC effect() {
  (FOR value IN [1] { RETURN STEP() })?
}
RETURN effect()`, func() (runtime.Value, error) {
					calls++

					return runtime.None, fmt.Errorf("boom")
				}, nil)

				if got, want := calls, 1; got != want {
					t.Fatalf("STEP calls = %d, want %d", got, want)
				}
			})

			t.Run("fallback", func(t *testing.T) {
				stepCalls := 0
				fallbackCalls := 0
				runDiscardedRecoveryQuery(t, level, `FUNC effect() {
  (FOR value IN [1] { RETURN STEP() }) ON ERROR RETURN FALLBACK()
}
RETURN effect()`, func() (runtime.Value, error) {
					stepCalls++

					return runtime.None, fmt.Errorf("boom")
				}, func() runtime.Value {
					fallbackCalls++

					return runtime.NewArrayWith(runtime.NewInt(1))
				})

				if got, want := stepCalls, 1; got != want {
					t.Fatalf("STEP calls = %d, want %d", got, want)
				}

				if got, want := fallbackCalls, 1; got != want {
					t.Fatalf("FALLBACK calls = %d, want %d", got, want)
				}
			})

			t.Run("retry succeeds", func(t *testing.T) {
				stepCalls := 0
				fallbackCalls := 0
				runDiscardedRecoveryQuery(t, level, `FUNC effect() {
  (FOR value IN [1] { RETURN STEP() }) ON ERROR RETRY 1 OR RETURN FALLBACK()
}
RETURN effect()`, func() (runtime.Value, error) {
					stepCalls++
					if stepCalls == 1 {
						return runtime.None, fmt.Errorf("boom")
					}

					return runtime.NewInt(1), nil
				}, func() runtime.Value {
					fallbackCalls++

					return runtime.NewArray(0)
				})

				if got, want := stepCalls, 2; got != want {
					t.Fatalf("STEP calls = %d, want %d", got, want)
				}

				if fallbackCalls != 0 {
					t.Fatalf("FALLBACK calls = %d, want 0", fallbackCalls)
				}
			})

			t.Run("retry falls back", func(t *testing.T) {
				stepCalls := 0
				fallbackCalls := 0
				runDiscardedRecoveryQuery(t, level, `FUNC effect() {
  (FOR value IN [1] { RETURN STEP() }) ON ERROR RETRY 1 OR RETURN FALLBACK()
}
RETURN effect()`, func() (runtime.Value, error) {
					stepCalls++

					return runtime.None, fmt.Errorf("boom")
				}, func() runtime.Value {
					fallbackCalls++

					return runtime.NewArrayWith(runtime.NewInt(1))
				})

				if got, want := stepCalls, 2; got != want {
					t.Fatalf("STEP calls = %d, want %d", got, want)
				}

				if got, want := fallbackCalls, 1; got != want {
					t.Fatalf("FALLBACK calls = %d, want %d", got, want)
				}
			})
		})
	}
}

func TestDiscardedForClosesReturnedResources(t *testing.T) {
	for _, level := range []compiler.OptimizationLevel{compiler.O0, compiler.O1} {
		t.Run(optimizationName(level), func(t *testing.T) {
			var closeCount atomic.Int32
			program, err := spec.Compile(`FOR value IN 1..3 { RETURN RESOURCE(value) }`, level)
			if err != nil {
				t.Fatalf("compile query: %v", err)
			}

			out, err := spec.Run(program, vm.WithFunction("RESOURCE", func(_ context.Context, args ...runtime.Value) (runtime.Value, error) {
				return newBodyCompletionResource(uint64(args[0].(runtime.Int)), &closeCount), nil
			}))
			if err != nil {
				t.Fatalf("run query: %v", err)
			}

			if got, want := string(out), "null"; got != want {
				t.Fatalf("result = %s, want %s", got, want)
			}

			if got, want := closeCount.Load(), int32(3); got != want {
				t.Fatalf("resource close count = %d, want %d", got, want)
			}
		})
	}
}

func runDiscardedRecoveryQuery(
	t *testing.T,
	level compiler.OptimizationLevel,
	query string,
	step func() (runtime.Value, error),
	fallback func() runtime.Value,
) {
	t.Helper()

	program, err := spec.Compile(query, level)
	if err != nil {
		t.Fatalf("compile query: %v", err)
	}

	opts := []vm.EnvironmentOption{
		vm.WithFunction("STEP", func(context.Context, ...runtime.Value) (runtime.Value, error) {
			return step()
		}),
	}

	if fallback != nil {
		opts = append(opts, vm.WithFunction("FALLBACK", func(context.Context, ...runtime.Value) (runtime.Value, error) {
			return fallback(), nil
		}))
	}

	out, err := spec.Run(program, opts...)
	if err != nil {
		t.Fatalf("run query: %v", err)
	}

	if got, want := string(out), "null"; got != want {
		t.Fatalf("result = %s, want %s", got, want)
	}
}

func TestReturnedForOperands(t *testing.T) {
	RunSpecs(t, []spec.Spec{
		Array(`RETURN FOR value IN [1, 2, 3] { RETURN value * 2 }`, []any{2, 4, 6}, "direct FOR return"),
		Array(`RETURN (FOR value IN [1, 2, 3] { RETURN value * 2 })`, []any{2, 4, 6}, "parenthesized FOR return"),
		Array(`RETURN DISTINCT FOR value IN [1, 1, 2] { RETURN value }`, []any{1, 2}, "direct distinct FOR return"),
		Array(`RETURN DISTINCT FOR DISTINCT IN [[1], [1], [2]] { RETURN (DISTINCT) }`, []any{[]any{1}, []any{2}}, "direct distinct FOR return keeps DISTINCT identifier escape"),
		Array(`
FUNC project(items) {
  FUNC double(item) => item * 2
  RETURN FOR item IN items {
    RETURN double(item)
  }
}
RETURN project([1, 2, 3])
`, []any{2, 4, 6}, "UDF direct FOR return preserves reachable nested calls"),
		Array(`
RETURN FOR outer IN [1, 2] {
  RETURN FOR inner IN [outer, outer + 10] {
    RETURN inner
  }
}
`, []any{[]any{1, 11}, []any{2, 12}}, "nested direct FOR return preserves nesting"),
	})
}
