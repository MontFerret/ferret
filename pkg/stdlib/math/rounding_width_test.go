package math_test

import (
	"context"
	"errors"
	stdmath "math"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/stdlib/math"
)

func TestRoundingIntegerWidth(t *testing.T) {
	functions := map[string]func(context.Context, runtime.Value) (runtime.Value, error){
		"floor": math.Floor, "ceil": math.Ceil, "round": math.Round,
	}
	for name, function := range functions {
		t.Run(name, func(t *testing.T) {
			for _, value := range []int64{stdmath.MinInt64, stdmath.MinInt32 - 1, stdmath.MinInt32, stdmath.MaxInt32, stdmath.MaxInt32 + 1, 1<<53 + 1, stdmath.MaxInt64} {
				input := runtime.NewInt64(value)
				got, err := function(t.Context(), input)
				if err != nil || got != input {
					t.Errorf("%s(Int(%d)) = %T(%v), %v", name, value, got, got, err)
				}
			}

			for _, value := range []float64{-0x1p63, stdmath.MinInt32 - 1, stdmath.MaxInt32 + 1, stdmath.Nextafter(0x1p63, 0)} {
				got, err := function(t.Context(), runtime.Float(value))
				if err != nil || got != runtime.NewInt64(int64(value)) {
					t.Errorf("%s(Float(%v)) = %T(%v), %v", name, value, got, got, err)
				}
			}

			for _, value := range []float64{stdmath.NaN(), stdmath.Inf(-1), stdmath.Inf(1), 0x1p63, stdmath.Nextafter(-0x1p63, stdmath.Inf(-1))} {
				got, err := function(t.Context(), runtime.Float(value))
				position, attributed, _ := runtime.InvalidArgumentDetails(err)
				if !errors.Is(err, runtime.ErrRange) || !attributed || position != 0 || got != runtime.None {
					t.Errorf("%s(%v) = %v, %v; want argument-0 range error", name, value, got, err)
				}
			}
		})
	}
}
