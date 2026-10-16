package runtime_test

import (
	"context"
	"errors"
	"fmt"
	"io"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func TestForEachIterRetainsJoinedEOFFailure(t *testing.T) {
	failure := errors.New("iterator failed")
	err := runtime.ForEachIter(t.Context(), failedIterator{err: errors.Join(io.EOF, failure)}, func(context.Context, runtime.Value, runtime.Value) (runtime.Boolean, error) {
		t.Fatal("failed yield visited")

		return false, nil
	})
	if !errors.Is(err, failure) || !errors.Is(err, io.EOF) {
		t.Fatalf("lost joined failure: %v", err)
	}
}

func TestForEachIterPreservesWrappedExhaustion(t *testing.T) {
	for _, err := range []error{io.EOF, fmt.Errorf("exhausted: %w", io.EOF), errors.Join(io.EOF, fmt.Errorf("wrapped: %w", io.EOF)), exhaustionError{}, fmt.Errorf("wrapped: %w", exhaustionError{})} {
		if got := runtime.ForEachIter(t.Context(), failedIterator{err: err}, func(context.Context, runtime.Value, runtime.Value) (runtime.Boolean, error) {
			t.Fatal("exhaustion visited")

			return false, nil
		}); got != nil {
			t.Fatalf("exhaustion %v: %v", err, got)
		}
	}
}
