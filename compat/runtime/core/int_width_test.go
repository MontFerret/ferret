package core_test

import (
	"math"
	"testing"

	"github.com/MontFerret/ferret/v2/compat/runtime/core"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func TestIntegerUnwrapPreservesWidth(t *testing.T) {
	for _, value := range []int64{math.MinInt64, math.MinInt32 - 1, 0, 42, math.MaxInt32 + 1, 1<<53 + 1, math.MaxInt64} {
		got := core.WrapValue(runtime.NewInt64(value)).Unwrap()
		if typed, ok := got.(int64); !ok || typed != value {
			t.Errorf("compat Unwrap = %T(%v), want int64(%d)", got, got, value)
		}
	}
}
