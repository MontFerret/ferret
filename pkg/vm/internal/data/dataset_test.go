package data_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/data"
)

func TestDistinctDataSetSeparatesHashCollisions(t *testing.T) {
	ctx := context.Background()
	set := data.NewDataSet(true)
	first := collisionValue{hash: 7, label: "first"}
	second := collisionValue{hash: 7, label: "second"}

	for _, value := range []runtime.Value{first, second, first, second} {
		if err := set.Append(ctx, value); err != nil {
			t.Fatalf("append %v: %v", value, err)
		}
	}

	length, err := set.Length(ctx)
	if err != nil {
		t.Fatalf("length: %v", err)
	}

	if length != 2 {
		t.Fatalf("expected 2 distinct values, got %d", length)
	}

	for idx, want := range []runtime.Value{first, second} {
		got, err := set.At(ctx, runtime.Int(idx))
		if err != nil {
			t.Fatalf("value at %d: %v", idx, err)
		}

		if runtime.CompareValues(got, want) != 0 {
			t.Fatalf("value at %d: expected %v, got %v", idx, want, got)
		}
	}
}
