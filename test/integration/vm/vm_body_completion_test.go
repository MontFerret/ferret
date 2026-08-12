package vm_test

import (
	"context"
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
