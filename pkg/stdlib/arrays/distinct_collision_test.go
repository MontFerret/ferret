package arrays_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/stdlib/arrays"
)

type distinctCollisionValue struct {
	label string
}

func (v distinctCollisionValue) String() string {
	return v.label
}

func (v distinctCollisionValue) Hash() uint64 {
	return 7
}

func (v distinctCollisionValue) Copy() runtime.Value {
	return v
}

func (v distinctCollisionValue) Compare(other runtime.Value) int {
	o, ok := other.(distinctCollisionValue)
	if !ok {
		return runtime.CompareTypes(v, other)
	}

	switch {
	case v.label < o.label:
		return -1
	case v.label > o.label:
		return 1
	default:
		return 0
	}
}

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

func TestUnionDistinctSeparatesHashCollisions(t *testing.T) {
	ctx := context.Background()
	first := distinctCollisionValue{label: "first"}
	second := distinctCollisionValue{label: "second"}

	result, err := arrays.UnionDistinct(
		ctx,
		runtime.NewArrayWith(first),
		runtime.NewArrayWith(second, first),
	)
	if err != nil {
		t.Fatalf("UnionDistinct: %v", err)
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

		if runtime.CompareValues(got, want) != 0 {
			t.Fatalf("value at %d: expected %v, got %v", idx, want, got)
		}
	}
}
