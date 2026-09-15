package math_test

import (
	"errors"
	stdmath "math"
	"math/rand/v2"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/rnd"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/stdlib/math"
)

func TestLegacyRandContracts(t *testing.T) {
	for _, test := range []struct {
		name     string
		args     []runtime.Value
		max, min float64
	}{
		{"unit", nil, 0, 0},
		{"one positive", []runtime.Value{runtime.Int(8)}, 16, 4},
		{"one negative", []runtime.Value{runtime.Int(-8)}, -16, -4},
		{"one zero", []runtime.Value{runtime.Int(0)}, 0, 0},
		{"one string", []runtime.Value{runtime.String("8")}, 16, 4},
		{"one bool", []runtime.Value{runtime.True}, 2, 0.5},
		{"max first", []runtime.Value{runtime.Int(6), runtime.Int(1)}, 6, 1},
		{"reversed", []runtime.Value{runtime.Int(1), runtime.Int(6)}, 1, 6},
		{"fractional", []runtime.Value{runtime.Float(2.75), runtime.Float(-1.25)}, 2.75, -1.25},
		{"coerced pair", []runtime.Value{runtime.String("6"), runtime.True}, 6, 1},
		{"equal", []runtime.Value{runtime.Int(5), runtime.Int(5)}, 5, 5},
		{"infinite max", []runtime.Value{runtime.Float(stdmath.Inf(1)), runtime.Int(0)}, stdmath.Inf(1), 0},
		{"infinite min", []runtime.Value{runtime.Int(0), runtime.Float(stdmath.Inf(-1))}, 0, stdmath.Inf(-1)},
		{"nan", []runtime.Value{runtime.Float(stdmath.NaN()), runtime.Int(0)}, stdmath.NaN(), 0},
	} {
		t.Run(test.name, func(t *testing.T) {
			src := rnd.NewSeed(42)
			ctx := rnd.WithContext(t.Context(), src)
			control := rand.New(rand.NewPCG(42, 0))
			for range 32 {
				unit := control.Float64()
				want := unit
				if len(test.args) != 0 {
					want = stdmath.Floor(unit*(test.max-test.min+1)) + test.min
				}

				got, err := math.Rand(ctx, test.args...)
				floating, ok := got.(runtime.Float)
				if err != nil || !ok || float64(floating) != want && !(stdmath.IsNaN(float64(floating)) && stdmath.IsNaN(want)) {
					t.Fatalf("got %v, %v; want Float(%v)", got, err, want)
				}
			}

			if src.Float64() != control.Float64() {
				t.Fatal("legacy draws did not advance the Session source exactly once per call")
			}
		})
	}
}

func TestLegacyRandErrors(t *testing.T) {
	for _, args := range [][]runtime.Value{
		{runtime.String("invalid")}, {runtime.Int(5), runtime.String("invalid")},
		{runtime.Int(1), runtime.Int(2), runtime.Int(3)},
	} {
		src, control := rnd.NewSeed(42), rnd.NewSeed(42)
		if got, err := math.Rand(rnd.WithContext(t.Context(), src), args...); err == nil || got != runtime.None {
			t.Fatalf("invalid call returned %v, %v", got, err)
		}

		if src.Float64() != control.Float64() {
			t.Fatal("argument errors consumed randomness")
		}
	}

	for _, args := range [][]runtime.Value{nil, {runtime.Int(5)}, {runtime.Int(5), runtime.Int(5)}} {
		if got, err := math.Rand(t.Context(), args...); got != runtime.None || !errors.Is(err, runtime.ErrUnexpected) {
			t.Fatalf("missing source returned %v, %v", got, err)
		}
	}
}
