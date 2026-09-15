package math

import (
	"context"
	"errors"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func TestSortNumbersCompletesWithCanceledContext(t *testing.T) {
	for _, ascending := range []bool{false, true} {
		ctx, cancel := context.WithCancel(t.Context())
		defer cancel()

		cancel()
		values := make([]runtime.Value, 4096)
		for index := range values {
			values[index] = runtime.Int((index * 37) % len(values))
		}
		if err := sortNumbers(ctx, values, ascending); err != nil {
			t.Fatal(err)
		}

		for index, value := range values {
			want := index
			if !ascending {
				want = len(values) - 1 - index
			}

			if value != runtime.Int(want) {
				t.Fatalf("position %d = %v, want %d", index, value, want)
			}
		}

		if err := sortNumbers(ctx, nil, ascending); err != nil {
			t.Fatalf("empty sort error = %v", err)
		}
	}
}

func TestSortNumbersCompletesAfterComparisonCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	value := &sortErrorValue{Value: runtime.ZeroInt, cancel: cancel}
	if err := sortNumbers(ctx, []runtime.Value{value, value}, true); err != nil {
		t.Fatal(err)
	}

	if ctx.Err() != context.Canceled {
		t.Fatal("comparison did not cancel the context")
	}
}

func TestSortNumbersOperationError(t *testing.T) {
	for _, ascending := range []bool{false, true} {
		for _, cancelDuringComparison := range []bool{false, true} {
			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()

			sentinel := errors.New("comparison failed")
			value := &sortErrorValue{Value: runtime.ZeroInt, err: sentinel}
			if cancelDuringComparison {
				value.cancel = cancel
			}

			if err := sortNumbers(ctx, []runtime.Value{value, value}, ascending); err != sentinel {
				t.Fatalf("sort error = %v, want original comparison error", err)
			}
		}
	}
}
