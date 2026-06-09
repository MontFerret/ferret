package vm

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
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

func TestArrayDistinctSeparatesHashCollisions(t *testing.T) {
	ctx := context.Background()
	first := distinctCollisionValue{label: "first"}
	second := distinctCollisionValue{label: "second"}

	result, err := arrayDistinct(ctx, runtime.NewArrayWith(first, second, first))
	if err != nil {
		t.Fatalf("arrayDistinct: %v", err)
	}

	length, err := result.Length(ctx)
	if err != nil {
		t.Fatalf("result length: %v", err)
	}

	if length != 2 {
		t.Fatalf("expected 2 distinct values, got %d", length)
	}

	for idx, want := range []runtime.Value{first, second} {
		got, err := result.At(ctx, runtime.Int(idx))
		if err != nil {
			t.Fatalf("result at %d: %v", idx, err)
		}

		if runtime.CompareValues(got, want) != 0 {
			t.Fatalf("result at %d: expected %v, got %v", idx, want, got)
		}
	}
}
