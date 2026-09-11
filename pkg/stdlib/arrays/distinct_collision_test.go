package arrays_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/stdlib/arrays"
)

func TestUniqueSeparatesHashCollisions(t *testing.T) {
	ctx := context.Background()
	first := distinctCollisionValue{label: "first"}
	second := distinctCollisionValue{label: "second"}

	result, err := arrays.Unique(ctx, runtime.NewArrayWith(first, second, first))
	if err != nil {
		t.Fatalf("Unique: %v", err)
	}

	assertDistinctCollisionValues(t, ctx, result.(runtime.List), first, second)
}

func TestUnionSeparatesHashCollisions(t *testing.T) {
	ctx := context.Background()
	first := distinctCollisionValue{label: "first"}
	second := distinctCollisionValue{label: "second"}

	result, err := arrays.Union(
		ctx,
		runtime.NewArrayWith(first),
		runtime.NewArrayWith(second, first),
	)
	if err != nil {
		t.Fatalf("Union: %v", err)
	}

	assertDistinctCollisionValues(t, ctx, result.(runtime.List), first, second)
}

func assertDistinctCollisionValues(t *testing.T, ctx context.Context, list runtime.List, expected ...runtime.Value) {
	t.Helper()

	length, err := list.Length(ctx)
	if err != nil {
		t.Fatalf("length: %v", err)
	}

	if length != runtime.Int(len(expected)) {
		t.Fatalf("expected %d values, got %d", len(expected), length)
	}

	for idx, want := range expected {
		got, err := list.At(ctx, runtime.Int(idx))
		if err != nil {
			t.Fatalf("value at %d: %v", idx, err)
		}

		equal, err := runtime.EqualValues(ctx, got, want)
		if err != nil {
			t.Fatalf("compare value at %d: %v", idx, err)
		}
		if !equal {
			t.Fatalf("value at %d: expected %v, got %v", idx, want, got)
		}
	}
}
