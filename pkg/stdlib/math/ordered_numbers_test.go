package math

import (
	"context"
	"errors"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func TestSortNumbersCancellation(t *testing.T) {
	for _, ascending := range []bool{false, true} {
		ctx, cancel := context.WithCancel(t.Context())
		defer cancel()

		polling := &sortCancelContext{Context: ctx, cancel: cancel, remaining: 3}
		values := []runtime.Value{runtime.Int(4), runtime.Int(1), runtime.Int(3), runtime.Int(2)}
		if err := sortNumbers(polling, values, ascending); !errors.Is(err, context.Canceled) {
			t.Fatalf("sort error = %v, want cancellation", err)
		}

		if err := sortNumbers(ctx, nil, ascending); !errors.Is(err, context.Canceled) {
			t.Fatalf("empty sort error = %v, want cancellation", err)
		}
	}
}
