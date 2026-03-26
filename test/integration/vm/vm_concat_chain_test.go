package vm_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
	"github.com/MontFerret/ferret/v2/test/spec"
	"github.com/MontFerret/ferret/v2/test/spec/assert"
	. "github.com/MontFerret/ferret/v2/test/spec/exec"
)

func TestConcatChainPreservesLeftToRightEvaluationOrder(t *testing.T) {
	RunSpecFactory(t, func() []spec.Spec {
		calls := make([]string, 0, 3)

		return []spec.Spec{
			spec.NewSpec(`RETURN TRACE_CALL("a") + "-x-" + TRACE_CALL("b") + "-y-" + TRACE_CALL("c")`).
				Env(vm.WithFunction("TRACE_CALL", func(_ context.Context, args ...runtime.Value) (runtime.Value, error) {
					if len(args) != 1 {
						return runtime.None, fmt.Errorf("unexpected TRACE_CALL arg count: got %d, want 1", len(args))
					}

					value := runtime.ToString(args[0]).String()
					calls = append(calls, value)

					return runtime.NewString(value), nil
				})).
				Expect().Exec(assert.NewUnaryAssertion(func(actual any) error {
				if err := assert.Equal(actual, "a-x-b-y-c"); err != nil {
					return err
				}

				if got, want := strings.Join(calls, ","), "a,b,c"; got != want {
					return fmt.Errorf("unexpected call order: got %q, want %q", got, want)
				}

				return nil
			})),
		}
	})
}
