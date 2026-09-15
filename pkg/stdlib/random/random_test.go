package random

import (
	"context"
	"errors"
	"math"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/rnd"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func TestFloatIntervals(t *testing.T) {
	for _, bounds := range [][2]runtime.Value{
		{runtime.Int(1), runtime.Int(6)}, {runtime.Float(-7.5), runtime.Int(-2)},
		{runtime.Int(-10), runtime.Float(0.5)},
		{runtime.Float(-math.MaxFloat64), runtime.Float(math.MaxFloat64)},
		{runtime.Float(1), runtime.Float(math.Nextafter(1, 2))},
		{runtime.Float(0), runtime.Float(math.SmallestNonzeroFloat64)},
		{runtime.Float(-math.SmallestNonzeroFloat64), runtime.Float(0)},
		{runtime.Int(1<<53 + 1), runtime.Int(1<<53 + 4)},
		{runtime.Int(1 << 53), runtime.Int(1<<53 + 1)},
		{runtime.Int(-1<<53 - 1), runtime.Int(-1<<53 + 1)},
		{runtime.Int(math.MaxInt64 - 4096), runtime.Int(math.MaxInt64)},
		{runtime.Int(math.MinInt64), runtime.Int(math.MinInt64 + 4096)},
	} {
		ctx := rnd.WithContext(t.Context(), rnd.NewSeed(42))
		for range 256 {
			got, err := Float(ctx, bounds[:]...)
			if err != nil {
				t.Fatalf("%v: %v", bounds, err)
			}

			floating, ok := got.(runtime.Float)
			if !ok || math.IsNaN(float64(floating)) || math.IsInf(float64(floating), 0) {
				t.Fatalf("%v: invalid Float %v", bounds, got)
			}

			lower, lowerErr := runtime.CompareValues(ctx, got, bounds[0])
			upper, upperErr := runtime.CompareValues(ctx, got, bounds[1])
			if lowerErr != nil || upperErr != nil || lower == runtime.Less || upper != runtime.Less {
				t.Fatalf("%v: result %v violated original half-open interval", bounds, got)
			}
		}
	}
}

func TestRandomSharedSequenceAndEqualBounds(t *testing.T) {
	src, control := rnd.NewSeed(42), rnd.NewSeed(42)
	ctx := rnd.WithContext(t.Context(), src)
	for range 32 {
		got, err := Float(ctx)
		if want := runtime.Float(control.Float64()); err != nil || got != want {
			t.Fatalf("unit float = %v, %v; want %v", got, err, want)
		}

		got, err = Float(ctx, runtime.Int(5), runtime.Float(5))
		if err != nil || got != runtime.Float(5) {
			t.Fatalf("equal float = %v, %v", got, err)
		}

		got, err = Int(ctx, runtime.Int(math.MaxInt64), runtime.Int(math.MaxInt64))
		if err != nil || got != runtime.Int(math.MaxInt64) {
			t.Fatalf("equal int = %v, %v", got, err)
		}

		got, err = Int(ctx, runtime.Int(math.MinInt64), runtime.Int(math.MaxInt64))
		if want := runtime.Int(control.Int64(math.MinInt64, math.MaxInt64)); err != nil || got != want {
			t.Fatalf("full int = %v, %v; want %v", got, err, want)
		}

		got, err = Bool(ctx)
		if want := runtime.Boolean(control.Bool()); err != nil || got != want {
			t.Fatalf("bool = %v, %v; want %v", got, err, want)
		}
	}

	negativeZero := runtime.Float(math.Copysign(0, -1))
	got, err := Float(ctx, negativeZero, runtime.Int(0))
	if err != nil || math.Float64bits(float64(got.(runtime.Float))) != math.Float64bits(float64(negativeZero)) {
		t.Fatalf("equal bounds lost signed zero: %v, %v", got, err)
	}

	got, err = Float(ctx, runtime.Int(math.MaxInt64), runtime.Int(math.MaxInt64))
	if err != nil || got != runtime.Float(math.MaxInt64) {
		t.Fatalf("equal bounds must use ordinary Float conversion: %v, %v", got, err)
	}

	if src.Float64() != control.Float64() {
		t.Fatal("equal bounds consumed randomness")
	}
}

func TestRandomArgumentErrorsDoNotDraw(t *testing.T) {
	for _, test := range []struct {
		want     error
		call     func(context.Context) (runtime.Value, error)
		name     string
		position int
	}{
		{name: "float string", call: func(ctx context.Context) (runtime.Value, error) {
			return Float(ctx, runtime.String("1"), runtime.Int(2))
		}, want: runtime.ErrInvalidType, position: 0},
		{name: "float boolean", call: func(ctx context.Context) (runtime.Value, error) { return Float(ctx, runtime.Int(1), runtime.True) }, want: runtime.ErrInvalidType, position: 1},
		{name: "float none", call: func(ctx context.Context) (runtime.Value, error) { return Float(ctx, runtime.None, runtime.Int(2)) }, want: runtime.ErrInvalidType, position: 0},
		{name: "reversed", call: func(ctx context.Context) (runtime.Value, error) { return Float(ctx, runtime.Int(2), runtime.Int(1)) }, want: runtime.ErrInvalidArgument, position: 0},
		{name: "reversed large ints", call: func(ctx context.Context) (runtime.Value, error) {
			return Float(ctx, runtime.Int(1<<53+1), runtime.Int(1<<53))
		}, want: runtime.ErrInvalidArgument, position: 0},
		{name: "no float", call: func(ctx context.Context) (runtime.Value, error) {
			return Float(ctx, runtime.Int(1<<53+1), runtime.Int(1<<53+2))
		}, want: runtime.ErrInvalidArgument, position: 1},
		{name: "int float min", call: func(ctx context.Context) (runtime.Value, error) { return Int(ctx, runtime.Float(1), runtime.Int(2)) }, want: runtime.ErrInvalidType, position: 0},
		{name: "int float max", call: func(ctx context.Context) (runtime.Value, error) { return Int(ctx, runtime.Int(1), runtime.Float(2)) }, want: runtime.ErrInvalidType, position: 1},
		{name: "int string", call: func(ctx context.Context) (runtime.Value, error) { return Int(ctx, runtime.String("1"), runtime.Int(2)) }, want: runtime.ErrInvalidType, position: 0},
		{name: "int reversed", call: func(ctx context.Context) (runtime.Value, error) { return Int(ctx, runtime.Int(2), runtime.Int(1)) }, want: runtime.ErrInvalidArgument, position: 0},
	} {
		t.Run(test.name, func(t *testing.T) {
			src, control := rnd.NewSeed(42), rnd.NewSeed(42)
			got, err := test.call(rnd.WithContext(t.Context(), src))
			position, ok, _ := runtime.InvalidArgumentDetails(err)
			if got != runtime.None || !errors.Is(err, test.want) || !ok || position != test.position {
				t.Fatalf("result %v, error %v; want %v at %d", got, err, test.want, test.position)
			}

			if src.Float64() != control.Float64() {
				t.Fatal("validation consumed randomness")
			}
		})
	}

	for _, bad := range []float64{math.NaN(), math.Inf(-1), math.Inf(1)} {
		for pos := range 2 {
			bounds := []runtime.Value{runtime.Int(0), runtime.Int(1)}
			bounds[pos] = runtime.Float(bad)
			src, control := rnd.NewSeed(42), rnd.NewSeed(42)
			got, err := Float(rnd.WithContext(t.Context(), src), bounds...)
			position, ok, _ := runtime.InvalidArgumentDetails(err)
			if got != runtime.None || !errors.Is(err, runtime.ErrInvalidArgument) || !ok || position != pos {
				t.Fatalf("non-finite %v: got %v, %v", bounds, got, err)
			}

			if src.Float64() != control.Float64() {
				t.Fatal("non-finite validation consumed randomness")
			}
		}
	}

	for _, bad := range []float64{math.NaN(), math.Inf(-1), math.Inf(1)} {
		src, control := rnd.NewSeed(42), rnd.NewSeed(42)
		got, err := Float(rnd.WithContext(t.Context(), src), runtime.Float(bad), runtime.Float(bad))
		if got != runtime.None || !errors.Is(err, runtime.ErrInvalidArgument) {
			t.Fatalf("equal non-finite bounds returned %v, %v", got, err)
		}

		if src.Float64() != control.Float64() {
			t.Fatal("equal non-finite bounds consumed randomness")
		}
	}

	for _, args := range [][]runtime.Value{{runtime.Int(1)}, {runtime.Int(1), runtime.Int(2), runtime.Int(3)}} {
		if got, err := Float(t.Context(), args...); got != runtime.None || !errors.Is(err, runtime.ErrInvalidArgumentNumber) {
			t.Fatalf("invalid arity: %v, %v", got, err)
		}
	}
}

func TestRandomMissingSource(t *testing.T) {
	for _, call := range []func(context.Context) (runtime.Value, error){
		func(ctx context.Context) (runtime.Value, error) { return Float(ctx) },
		func(ctx context.Context) (runtime.Value, error) { return Float(ctx, runtime.Int(0), runtime.Int(1)) },
		func(ctx context.Context) (runtime.Value, error) { return Float(ctx, runtime.Int(5), runtime.Int(5)) },
		func(ctx context.Context) (runtime.Value, error) { return Int(ctx, runtime.Int(0), runtime.Int(1)) },
		func(ctx context.Context) (runtime.Value, error) { return Int(ctx, runtime.Int(5), runtime.Int(5)) },
		Bool,
	} {
		if got, err := call(t.Context()); got != runtime.None || !errors.Is(err, runtime.ErrUnexpected) {
			t.Fatalf("missing source: %v, %v", got, err)
		}
	}
}

func TestFloatInterpolationRounding(t *testing.T) {
	for _, bounds := range [][2]float64{
		{1, math.Nextafter(1, 2)}, {-math.MaxFloat64, math.MaxFloat64},
		{-math.MaxFloat64, -math.MaxFloat64 / 2}, {math.MaxFloat64 / 2, math.MaxFloat64},
		{0, math.SmallestNonzeroFloat64}, {-math.SmallestNonzeroFloat64, 0},
	} {
		for _, unit := range []float64{0, 0.5, math.Nextafter(1, 0)} {
			got := interpolateFloat(bounds[0], bounds[1], unit)
			if math.IsNaN(got) || got < bounds[0] || got >= bounds[1] {
				t.Fatalf("%v with %g: got %g", bounds, unit, got)
			}
		}
	}
}
