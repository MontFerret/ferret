package vm_test

import (
	"context"
	"strings"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

func TestConcatChainPreservesLeftToRightEvaluationOrder(t *testing.T) {
	forEachCompiledVM(t, `RETURN LOG("a") + "-x-" + LOG("b") + "-y-" + LOG("c")`, func(t *testing.T, instance *vm.VM) {
		calls := make([]string, 0, 3)
		env := mustNewEnvironment(t, vm.WithFunction("LOG", func(_ context.Context, args ...runtime.Value) (runtime.Value, error) {
			if len(args) != 1 {
				t.Fatalf("unexpected LOG arg count: got %d, want 1", len(args))
			}

			value := runtime.ToString(args[0]).String()
			calls = append(calls, value)

			return runtime.NewString(value), nil
		}))

		out, err := runVM(t, instance, env)
		if err != nil {
			t.Fatalf("unexpected runtime error: %v", err)
		}

		assertRuntimeValueEquals(t, out, runtime.NewString("a-x-b-y-c"))

		if got, want := strings.Join(calls, ","), "a,b,c"; got != want {
			t.Fatalf("unexpected call order: got %q, want %q", got, want)
		}
	})
}
