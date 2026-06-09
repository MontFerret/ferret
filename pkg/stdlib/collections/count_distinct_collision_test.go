package collections_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/stdlib/collections"
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

func TestCountDistinctSeparatesHashCollisions(t *testing.T) {
	first := distinctCollisionValue{label: "first"}
	second := distinctCollisionValue{label: "second"}

	result, err := collections.CountDistinct(
		context.Background(),
		runtime.NewArrayWith(first, second, first, second),
	)
	if err != nil {
		t.Fatalf("CountDistinct: %v", err)
	}

	if result != runtime.NewInt(2) {
		t.Fatalf("expected 2 distinct values, got %v", result)
	}
}
