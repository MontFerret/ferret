package vm_test

import (
	"context"
	"fmt"
	"math"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/compiler"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/source"
	"github.com/MontFerret/ferret/v2/pkg/vm"
	"github.com/MontFerret/ferret/v2/test/spec"
	. "github.com/MontFerret/ferret/v2/test/spec/exec"
)

func TestRetryCountIntegerWidth(t *testing.T) {
	for _, level := range []compiler.OptimizationLevel{compiler.None, compiler.Basic, compiler.Full} {
		for _, count := range []int64{math.MaxInt32, math.MaxInt32 + 1, 1<<53 + 1, math.MaxInt64} {
			t.Run(fmt.Sprintf("%s/%d", level, count), func(t *testing.T) {
				calls := 0
				RunSpecsWith(t, "retry",
					mustNewCompiler(t, compiler.WithOptimizationLevel(level)),
					[]spec.Spec{S(fmt.Sprintf("RETURN STEP() ON ERROR RETRY %d", count), 99)},
					vm.WithFunction("STEP", func(context.Context, ...runtime.Value) (runtime.Value, error) {
						calls++
						if calls == 1 {
							return runtime.None, fmt.Errorf("retry once")
						}
						return runtime.Int(99), nil
					}),
				)
				if calls != 2 {
					t.Fatalf("calls = %d, want 2", calls)
				}
			})
		}

		compiler := mustNewCompiler(t, compiler.WithOptimizationLevel(level))
		_, err := compiler.Compile(t.Context(), source.NewAnonymous("RETURN 42 ON ERROR RETRY 9223372036854775808"))
		if err == nil {
			t.Fatal("out-of-range retry count compiled")
		}
	}
}
