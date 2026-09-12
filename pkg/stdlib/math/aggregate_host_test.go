package math_test

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"
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
			source := &aggregateList{values: values}
			got, err := fn.call(ctx, source)
			if err != nil {
				t.Fatal(err)
			}

			assertMathResult(t, got, fn.want[0])
			if source.context != ctx || !slices.Equal(values, before) {
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
	sentinel := errors.New("host read failed")
	for _, fn := range aggregateCases() {
		for _, after := range []int{0, 1, 3} {
			t.Run(fmt.Sprintf("%s/after-%d", fn.name, after), func(t *testing.T) {
				source := &aggregateList{
					values: []runtime.Value{runtime.Int(3), runtime.Int(1), runtime.Int(2)},
					visit: func(_ context.Context, _, index int) error {
						if index == after {
							return fmt.Errorf("read: %w", sentinel)
						}

						return nil
					},
				}
				_, err := fn.call(t.Context(), source)
				if !errors.Is(err, sentinel) {
					t.Fatalf("error = %v, want host error", err)
				}
			})
		}

		t.Run(fn.name+"/empty-error", func(t *testing.T) {
			source := &aggregateList{visit: func(context.Context, int, int) error { return sentinel }}
			_, err := fn.call(t.Context(), source)
			if !errors.Is(err, sentinel) {
				t.Fatalf("error = %v, want error before empty result", err)
			}
		})

		if !strings.HasPrefix(fn.name, "variance") && !strings.HasPrefix(fn.name, "stddev") {
			continue
		}

		t.Run(fn.name+"/second-pass", func(t *testing.T) {
			source := &aggregateList{
				values: []runtime.Value{runtime.Int(3), runtime.Int(1), runtime.Int(2)},
				visit: func(_ context.Context, pass, index int) error {
					if pass == 2 && index == 1 {
						return sentinel
					}

					return nil
				},
			}
			_, err := fn.call(t.Context(), source)
			if !errors.Is(err, sentinel) {
				t.Fatalf("second-pass error = %v", err)
			}
		})
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

		for _, after := range []int{1, 3} {
			t.Run(fmt.Sprintf("%s/cancel-after-%d", fn.name, after), func(t *testing.T) {
				ctx, cancel := context.WithCancel(t.Context())
				defer cancel()

				source := &aggregateList{
					values: []runtime.Value{runtime.Int(3), runtime.Int(1), runtime.Int(2)},
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

		t.Run(fn.name+"/error-and-cancellation", func(t *testing.T) {
			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()

			sentinel := errors.New("host read failed")
			source := &aggregateList{visit: func(context.Context, int, int) error {
				cancel()

				return sentinel
			}}
			_, err := fn.call(ctx, source)
			if !errors.Is(err, context.Canceled) || !errors.Is(err, sentinel) {
				t.Fatalf("error = %v, want both host error and cancellation", err)
			}
		})
	}
}
