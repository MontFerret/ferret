package math_test

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/stdlib/math"
)

func TestSumPropagatesIterationError(t *testing.T) {
	sentinel := errors.New("host read failed")
	source := &aggregateList{
		List:   runtime.NewArrayWith(runtime.Int(1)),
		values: []runtime.Value{runtime.Int(1)},
		visit: func(_ context.Context, _, index int) error {
			if index == 1 {
				return sentinel
			}

			return nil
		},
	}
	_, err := math.Sum(t.Context(), source)
	if !errors.Is(err, sentinel) {
		t.Fatalf("error = %v, want iteration error after partial sum", err)
	}
}

func TestMathAggregatesHostLists(t *testing.T) {
	for _, fn := range aggregateCases() {
		t.Run(fn.name, func(t *testing.T) {
			values := []runtime.Value{runtime.Int(5), runtime.Int(2), runtime.Int(9), runtime.Int(2)}
			before := slices.Clone(values)
			ctx := t.Context()
			source := &aggregateList{
				values: values,
				visit: func(_ context.Context, pass, _ int) error {
					if pass > 1 {
						return errors.New("source cannot be traversed twice")
					}

					return nil
				},
			}
			got, err := fn.call(ctx, source)
			if err != nil {
				t.Fatal(err)
			}

			assertMathResult(t, got, fn.want[0])
			if source.calls != 1 || source.context != ctx || !slices.Equal(values, before) {
				t.Fatal("source context or contents changed")
			}

			empty := &aggregateList{}
			got, err = fn.call(ctx, empty)
			if err != nil {
				t.Fatal(err)
			}

			assertMathResult(t, got, fn.want[5])
		})
	}
}

func TestMathAggregatesIterationErrors(t *testing.T) {
	for _, fn := range aggregateCases() {
		for errorIndex, sentinel := range []error{errors.New("host read failed"), fmt.Errorf("host decode: %w", runtime.ErrInvalidType)} {
			for _, values := range [][]runtime.Value{nil, {runtime.None, runtime.Int(3), runtime.String("2"), runtime.Int(1), runtime.Int(2)}} {
				for after := 0; after <= len(values); after++ {
					t.Run(fmt.Sprintf("%s/error-%d/size-%d/after-%d", fn.name, errorIndex, len(values), after), func(t *testing.T) {
						source := &aggregateList{
							values: values,
							visit: func(_ context.Context, _, index int) error {
								if index == after {
									return sentinel
								}

								return nil
							},
						}
						_, err := fn.call(t.Context(), source)
						if err != sentinel {
							t.Fatalf("error = %v, want original host error %v", err, sentinel)
						}
					})
				}
			}
		}
	}
}

func TestMathAggregatesCancellation(t *testing.T) {
	for _, fn := range aggregateCases() {
		for _, empty := range []bool{false, true} {
			t.Run(fmt.Sprintf("%s/pre-canceled/empty=%t", fn.name, empty), func(t *testing.T) {
				ctx, cancel := context.WithCancel(t.Context())
				cancel()
				source := runtime.NewArray(0)
				if !empty {
					source = runtime.NewArrayWith(runtime.Int(1))
				}

				_, err := fn.call(ctx, source)
				if !errors.Is(err, context.Canceled) {
					t.Fatalf("error = %v, want cancellation", err)
				}
			})
		}

		for inputIndex, values := range [][]runtime.Value{
			nil,
			{runtime.Int(3), runtime.Int(1), runtime.Int(2)},
			{runtime.None, runtime.String("2"), runtime.True},
			{runtime.Int(1), runtime.None, runtime.Int(3)},
		} {
			for after := 0; after <= len(values); after++ {
				t.Run(fmt.Sprintf("%s/input-%d/cancel-after-%d", fn.name, inputIndex, after), func(t *testing.T) {
					ctx, cancel := context.WithCancel(t.Context())
					defer cancel()

					source := &aggregateList{
						values: values,
						visit: func(_ context.Context, _, index int) error {
							if index == after {
								cancel()
							}

							return nil
						},
					}
					_, err := fn.call(ctx, source)
					if !errors.Is(err, context.Canceled) {
						t.Fatalf("error = %v, want cancellation", err)
					}
				})
			}
		}

		t.Run(fn.name+"/error-and-cancellation", func(t *testing.T) {
			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()

			sentinel := errors.New("host read failed")
			source := &aggregateList{visit: func(context.Context, int, int) error {
				cancel()

				return sentinel
			}}
			_, err := fn.call(ctx, source)
			if err != sentinel || errors.Is(err, context.Canceled) {
				t.Fatalf("error = %v, want causal host error without joined cancellation", err)
			}
		})
	}
}
