package math_test

import (
	"context"
	"errors"
	"fmt"
	stdmath "math"
	"strings"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/stdlib/math"
)

type aggregateCase struct {
	name string
	call func(context.Context, runtime.Value) (runtime.Value, error)
	want []runtime.Value
}

func TestMathAggregates(t *testing.T) {
	inputs := []struct {
		name   string
		values []runtime.Value
	}{
		{"integers", []runtime.Value{runtime.Int(5), runtime.Int(2), runtime.Int(9), runtime.Int(2)}},
		{"floats", []runtime.Value{runtime.Float(1.5), runtime.Float(2.5), runtime.Float(3.5)}},
		{"mixed", []runtime.Value{runtime.Int(-3), runtime.Float(1), runtime.Int(2)}},
		{"negative", []runtime.Value{runtime.Int(-3), runtime.Int(-5), runtime.Int(-1)}},
		{"singleton", []runtime.Value{runtime.Int(7)}},
		{"empty", nil},
	}
	for _, fn := range aggregateCases() {
		for idx, input := range inputs {
			t.Run(fn.name+"/"+input.name, func(t *testing.T) {
				source := runtime.NewArrayWith(input.values...)
				before := source.String()
				got, err := fn.call(t.Context(), source)
				if err != nil {
					t.Fatal(err)
				}

				assertMathResult(t, got, fn.want[idx])
				if source.String() != before {
					t.Fatalf("input changed: %s -> %s", before, source)
				}
			})
		}
	}
}

func TestMathAggregatesRejectInvalidValues(t *testing.T) {
	invalid := []runtime.Value{runtime.String("2"), runtime.None, runtime.False, runtime.NewArray(0), runtime.NewObject(), runtime.Duration(1)}
	for _, fn := range aggregateCases() {
		for _, value := range invalid {
			for index := 0; index < 3; index++ {
				t.Run(fmt.Sprintf("%s/%T/%d", fn.name, value, index), func(t *testing.T) {
					values := []runtime.Value{runtime.Int(1), runtime.Int(2), runtime.Int(3)}
					values[index] = value
					source := runtime.NewArrayWith(values...)
					before := source.String()
					_, err := fn.call(t.Context(), source)
					if !errors.Is(err, runtime.ErrInvalidType) || !strings.Contains(err.Error(), fmt.Sprintf("at index %d", index)) {
						t.Fatalf("error = %v, want indexed type error", err)
					}

					if source.String() != before {
						t.Fatalf("invalid input changed: %s -> %s", before, source)
					}
				})
			}
		}

		t.Run(fn.name+"/non-list", func(t *testing.T) {
			_, err := fn.call(t.Context(), runtime.Int(1))
			if !errors.Is(err, runtime.ErrInvalidType) {
				t.Fatalf("error = %v, want type error", err)
			}
		})

		t.Run(fn.name+"/invalid-singleton", func(t *testing.T) {
			_, err := fn.call(t.Context(), runtime.NewArrayWith(runtime.None))
			if !errors.Is(err, runtime.ErrInvalidType) {
				t.Fatalf("error = %v, want type error before insufficient-input result", err)
			}
		})
	}
}

func TestMathAggregatesAcceptNonFiniteNumbers(t *testing.T) {
	for _, fn := range aggregateCases() {
		for _, value := range []runtime.Float{runtime.NaN(), runtime.Float(stdmath.Inf(1)), runtime.Float(stdmath.Inf(-1))} {
			t.Run(fn.name+"/"+value.String(), func(t *testing.T) {
				_, err := fn.call(t.Context(), runtime.NewArrayWith(value, runtime.Int(1)))
				if err != nil {
					t.Fatalf("non-finite numeric value rejected: %v", err)
				}
			})
		}
	}
}

func TestMedianMiddleElements(t *testing.T) {
	for _, test := range []struct {
		want   runtime.Value
		values []runtime.Value
	}{
		{values: []runtime.Value{runtime.Int(1), runtime.Int(2), runtime.Int(3)}, want: runtime.Int(2)},
		{values: []runtime.Value{runtime.Int(1), runtime.Int(2), runtime.Int(3), runtime.Int(4)}, want: runtime.Float(2.5)},
		{values: []runtime.Value{runtime.Int(100), runtime.Int(3), runtime.Int(1), runtime.Int(2)}, want: runtime.Float(2.5)},
		{values: []runtime.Value{runtime.Float(100), runtime.Int(-5), runtime.Float(3), runtime.Int(1)}, want: runtime.Float(2)},
		{values: []runtime.Value{runtime.Int(9007199254740994), runtime.Int(9007199254740993), runtime.Int(9007199254740992)}, want: runtime.Int(9007199254740993)},
	} {
		source := runtime.NewArrayWith(test.values...)
		t.Run(source.String(), func(t *testing.T) {
			before := source.String()
			got, err := math.Median(t.Context(), source)
			if err != nil {
				t.Fatal(err)
			}

			assertMathResult(t, got, test.want)
			if source.String() != before {
				t.Fatalf("input changed: %s -> %s", before, source)
			}
		})
	}
}

func aggregateCases() []aggregateCase {
	return []aggregateCase{
		{"sum", math.Sum, []runtime.Value{runtime.Float(18), runtime.Float(7.5), runtime.Float(0), runtime.Float(-9), runtime.Float(7), runtime.ZeroInt}},
		{"average", math.Average, []runtime.Value{runtime.Float(4.5), runtime.Float(2.5), runtime.Float(0), runtime.Float(-3), runtime.Float(7), runtime.NaN()}},
		{"min", math.Min, []runtime.Value{runtime.Float(2), runtime.Float(1.5), runtime.Float(-3), runtime.Float(-5), runtime.Float(7), runtime.None}},
		{"max", math.Max, []runtime.Value{runtime.Float(9), runtime.Float(3.5), runtime.Float(2), runtime.Float(-1), runtime.Float(7), runtime.None}},
		{"median", math.Median, []runtime.Value{runtime.Float(3.5), runtime.Float(2.5), runtime.Float(1), runtime.Int(-3), runtime.Int(7), runtime.NaN()}},
		{"variance_population", math.PopulationVariance, []runtime.Value{runtime.Float(8.25), runtime.Float(2.0 / 3), runtime.Float(14.0 / 3), runtime.Float(8.0 / 3), runtime.Float(0), runtime.NaN()}},
		{"variance_sample", math.SampleVariance, []runtime.Value{runtime.Float(11), runtime.Float(1), runtime.Float(7), runtime.Float(4), runtime.NaN(), runtime.NaN()}},
		{"stddev_population", math.StandardDeviationPopulation, []runtime.Value{runtime.Float(stdmath.Sqrt(8.25)), runtime.Float(stdmath.Sqrt(2.0 / 3)), runtime.Float(stdmath.Sqrt(14.0 / 3)), runtime.Float(stdmath.Sqrt(8.0 / 3)), runtime.Float(0), runtime.NaN()}},
		{"stddev_sample", math.StandardDeviationSample, []runtime.Value{runtime.Float(stdmath.Sqrt(11)), runtime.Float(1), runtime.Float(stdmath.Sqrt(7)), runtime.Float(2), runtime.NaN(), runtime.NaN()}},
		{"percentile_rank", func(ctx context.Context, value runtime.Value) (runtime.Value, error) {
			return math.Percentile(ctx, value, runtime.Int(50))
		}, []runtime.Value{runtime.Int(2), runtime.Float(2.5), runtime.Float(1), runtime.Int(-3), runtime.Int(7), runtime.NaN()}},
		{"percentile_interpolation", func(ctx context.Context, value runtime.Value) (runtime.Value, error) {
			return math.Percentile(ctx, value, runtime.Int(50), runtime.String("interpolation"))
		}, []runtime.Value{runtime.Float(3.5), runtime.Float(2.5), runtime.Float(1), runtime.Int(-3), runtime.Int(7), runtime.NaN()}},
	}
}

func assertMathResult(t *testing.T, got, want runtime.Value) {
	t.Helper()

	if expected, ok := want.(runtime.Float); ok {
		actual, ok := got.(runtime.Float)
		if !ok {
			t.Fatalf("result = %v (%T), want Float %v", got, got, want)
		}

		if runtime.IsNaN(expected) {
			if !runtime.IsNaN(actual) {
				t.Fatalf("result = %v, want NaN", got)
			}

			return
		}

		if runtime.IsNaN(actual) || stdmath.Abs(float64(actual-expected)) > 1e-12*stdmath.Max(1, stdmath.Abs(float64(expected))) {
			t.Fatalf("result = %v, want %v", got, want)
		}

		return
	}

	if got != want {
		t.Fatalf("result = %v (%T), want %v (%T)", got, got, want, want)
	}
}
