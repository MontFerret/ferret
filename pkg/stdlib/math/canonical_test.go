package math_test

import (
	"context"
	"errors"
	"fmt"
	stdmath "math"
	"slices"
	"strings"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/stdlib/math"
)

func TestCanonicalMathRejectsNonNumbers(t *testing.T) {
	for _, fn := range canonicalAggregateCases(t) {
		for _, invalid := range []runtime.Value{runtime.None, runtime.String("2"), runtime.True, runtime.False, runtime.NewArray(0), runtime.NewObject(), runtime.Duration(1)} {
			for index := 0; index <= 3; index++ {
				for _, host := range []bool{false, true} {
					t.Run(fmt.Sprintf("%s/%T/%d/host=%t", fn.name, invalid, index, host), func(t *testing.T) {
						values := slices.Insert([]runtime.Value{runtime.Int(1), runtime.Float(2), runtime.Int(3)}, index, invalid)
						before := slices.Clone(values)
						var source runtime.List = runtime.NewArrayWith(values...)
						if host {
							source = &aggregateList{values: values}
						}

						_, err := fn.call(t.Context(), source)
						if !errors.Is(err, runtime.ErrInvalidType) || !strings.Contains(err.Error(), fmt.Sprintf("index %d", index)) || !strings.Contains(err.Error(), "position 1") {
							t.Fatalf("error = %v, want indexed argument type error", err)
						}

						if !slices.Equal(values, before) {
							t.Fatal("source changed")
						}
					})
				}
			}
		}
	}
}

func TestCanonicalMathHostLists(t *testing.T) {
	for _, fn := range canonicalAggregateCases(t) {
		t.Run(fn.name, func(t *testing.T) {
			values := []runtime.Value{runtime.Int(5), runtime.Int(2), runtime.Int(9), runtime.Int(2)}
			source := &aggregateList{values: values, visit: func(_ context.Context, pass, _ int) error {
				if pass > 1 {
					return errors.New("source cannot be traversed twice")
				}

				return nil
			}}
			got, err := fn.call(t.Context(), source)
			if err != nil {
				t.Fatal(err)
			}

			assertMathResult(t, got, fn.want[0])
			if source.calls != 1 {
				t.Fatalf("traversals = %d", source.calls)
			}

			got, err = fn.call(t.Context(), &aggregateList{})
			if err != nil {
				t.Fatal(err)
			}

			assertMathResult(t, got, fn.want[5])
			for after := 0; after <= len(values); after++ {
				for _, cancelWithError := range []bool{false, true} {
					ctx, cancel := context.WithCancel(t.Context())
					sentinel := errors.New("host read failure")
					failed := &aggregateList{values: values, visit: func(_ context.Context, _, index int) error {
						if index == after {
							if cancelWithError {
								cancel()
							}

							return sentinel
						}

						return nil
					}}
					got, err = fn.call(ctx, failed)
					cancel()
					if err != sentinel || got != runtime.None {
						t.Fatalf("result = %v, error = %v, want original error without partial result", got, err)
					}
				}

				ctx, cancel := context.WithCancel(t.Context())
				canceled := &aggregateList{values: values, visit: func(_ context.Context, _, index int) error {
					if index == after {
						cancel()
					}

					return nil
				}}
				got, err = fn.call(ctx, canceled)
				cancel()
				if err != nil {
					t.Fatal(err)
				}

				assertMathResult(t, got, fn.want[0])
			}

			ctx, cancel := context.WithCancel(t.Context())
			cancel()
			got, err = fn.call(ctx, &aggregateList{})
			if err != nil {
				t.Fatalf("empty canceled list: %v", err)
			}

			assertMathResult(t, got, fn.want[5])
		})
	}
}

func TestCanonicalPercentile(t *testing.T) {
	fn, _ := mathFunctions(t).A2().Get("math::percentile")
	for _, p := range []runtime.Value{runtime.Int(0), runtime.Int(25), runtime.Int(50), runtime.Int(75), runtime.Int(95), runtime.Float(99.9), runtime.Int(100), runtime.Float(0), runtime.Float(100)} {
		source := runtime.NewArrayWith(runtime.Int(400), runtime.Int(300), runtime.Int(200), runtime.Int(100), runtime.Int(0))
		before := source.String()
		got, err := fn(t.Context(), source, p)
		if err != nil {
			t.Fatal(err)
		}

		percent, _ := runtime.ToFloat(t.Context(), p)
		want := runtime.Value(runtime.Float(percent * 4))
		if percent == 0 || percent == 25 || percent == 50 || percent == 75 || percent == 100 {
			want = runtime.Int(percent * 4)
		}

		assertMathResult(t, got, want)
		if source.String() != before {
			t.Fatal("percentile changed source")
		}
	}

	for _, p := range []runtime.Value{runtime.Int(-1), runtime.Int(101), runtime.Float(-0.01), runtime.Float(100.01), runtime.NaN(), runtime.Float(stdmath.Inf(1)), runtime.Float(stdmath.Inf(-1)), runtime.String("50"), runtime.None, runtime.True} {
		for _, source := range []*runtime.Array{runtime.NewArray(0), runtime.NewArrayWith(runtime.Int(1))} {
			got, err := fn(t.Context(), source, p)
			expected := runtime.ErrInvalidArgument
			if !runtime.IsNumber(p) {
				expected = runtime.ErrInvalidType
			}

			if got != runtime.None || !errors.Is(err, expected) {
				t.Fatalf("p=%v: result=%v error=%v, want %v", p, got, err, expected)
			}
		}
	}

	for _, p := range []runtime.Value{runtime.Int(0), runtime.Float(99.9), runtime.Int(100)} {
		got, err := fn(t.Context(), runtime.NewArrayWith(runtime.Int(9007199254740993)), p)
		if err != nil {
			t.Fatal(err)
		}

		assertMathResult(t, got, runtime.Int(9007199254740993))
	}
}

func TestCanonicalVarianceLargeOffset(t *testing.T) {
	for _, test := range []struct {
		name string
		want runtime.Float
	}{
		{"variance", 2}, {"variance_sample", 2.5}, {"stddev", runtime.Float(stdmath.Sqrt(2))}, {"stddev_sample", runtime.Float(stdmath.Sqrt(2.5))},
	} {
		fn, _ := mathFunctions(t).A1().Get("math::" + test.name)
		source := &aggregateList{values: []runtime.Value{runtime.Int(1000000000001), runtime.Float(1000000000002), runtime.Int(1000000000003), runtime.Float(1000000000004), runtime.Int(1000000000005)}}
		got, err := fn(t.Context(), source)
		if err != nil {
			t.Fatal(err)
		}

		assertMathResult(t, got, test.want)
	}
}

func mathFunctions(t testing.TB) *runtime.Functions {
	t.Helper()

	lib := runtime.NewLibrary()
	math.RegisterLib(lib)
	functions, err := lib.Build()
	if err != nil {
		t.Fatal(err)
	}

	return functions
}

func canonicalAggregateCases(t testing.TB) []aggregateCase {
	t.Helper()

	functions := mathFunctions(t)
	cases := aggregateCases()
	result := make([]aggregateCase, 0, len(cases)-1)
	for _, test := range cases {
		switch test.name {
		case "percentile_rank":
			continue
		case "percentile_interpolation":
			fn, _ := functions.A2().Get("math::percentile")
			test.name = "percentile"
			test.call = func(ctx context.Context, arg runtime.Value) (runtime.Value, error) {
				return fn(ctx, arg, runtime.Int(50))
			}

		default:
			switch test.name {
			case "average":
				test.name = "mean"
				test.want[5] = runtime.NaN()
			case "median":
				test.want[5] = runtime.NaN()
			case "variance_population":
				test.name = "variance"
			case "stddev_population":
				test.name = "stddev"
			}

			fn, ok := functions.A1().Get("math::" + test.name)
			if !ok {
				t.Fatalf("missing canonical function %s", test.name)
			}

			test.call = fn
		}

		test.name = "math::" + test.name
		result = append(result, test)
	}

	return result
}
