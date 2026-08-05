package arrays_test

import (
	"context"
	"errors"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/stdlib/arrays"
)

type failingEqualityValue struct {
	err error
}

func (v failingEqualityValue) String() string {
	return "failing"
}

func (v failingEqualityValue) Hash() uint64 {
	return 11
}

func (v failingEqualityValue) Copy() runtime.Value {
	return v
}

func (v failingEqualityValue) Equal(context.Context, runtime.Value) (bool, error) {
	return false, v.err
}

func TestSetOperationsVerifyHashCollisions(t *testing.T) {
	ctx := context.Background()
	first := distinctCollisionValue{label: "first"}
	second := distinctCollisionValue{label: "second"}

	intersection, err := arrays.Intersection(
		ctx,
		runtime.NewArrayWith(first, first, second),
		runtime.NewArrayWith(second, first),
	)
	if err != nil {
		t.Fatalf("Intersection: %v", err)
	}
	assertListContainsExactly(t, ctx, intersection.(runtime.List), first, second)

	outersection, err := arrays.Outersection(
		ctx,
		runtime.NewArrayWith(first, first),
		runtime.NewArrayWith(second),
	)
	if err != nil {
		t.Fatalf("Outersection: %v", err)
	}
	assertListContainsExactly(t, ctx, outersection.(runtime.List), first, second)

	minus, err := arrays.Minus(
		ctx,
		runtime.NewArrayWith(first, second, first),
		runtime.NewArrayWith(first),
	)
	if err != nil {
		t.Fatalf("Minus: %v", err)
	}
	assertListContainsExactly(t, ctx, minus.(runtime.List), second)

	remaining, err := arrays.RemoveValues(
		ctx,
		runtime.NewArrayWith(first, second),
		runtime.NewArrayWith(first),
	)
	if err != nil {
		t.Fatalf("RemoveValues: %v", err)
	}
	assertListContainsExactly(t, ctx, remaining.(runtime.List), second)
}

func TestSetOperationsUseNumericEqualityAcrossRepresentations(t *testing.T) {
	ctx := context.Background()
	integers := runtime.NewArrayWith(runtime.NewInt(1))
	floats := runtime.NewArrayWith(runtime.NewFloat(1))

	intersection, err := arrays.Intersection(ctx, integers, floats)
	if err != nil {
		t.Fatalf("Intersection: %v", err)
	}
	assertListContainsExactly(t, ctx, intersection.(runtime.List), runtime.NewInt(1))

	outersection, err := arrays.Outersection(ctx, integers, floats)
	if err != nil {
		t.Fatalf("Outersection: %v", err)
	}
	assertListContainsExactly(t, ctx, outersection.(runtime.List))

	minus, err := arrays.Minus(ctx, integers, floats)
	if err != nil {
		t.Fatalf("Minus: %v", err)
	}
	assertListContainsExactly(t, ctx, minus.(runtime.List))

	remaining, err := arrays.RemoveValues(ctx, integers, floats)
	if err != nil {
		t.Fatalf("RemoveValues: %v", err)
	}
	assertListContainsExactly(t, ctx, remaining.(runtime.List))

	union, err := arrays.UnionDistinct(ctx, integers, floats)
	if err != nil {
		t.Fatalf("UnionDistinct: %v", err)
	}
	assertListContainsExactly(t, ctx, union.(runtime.List), runtime.NewInt(1))
}

func TestHashSetOperationsPropagateEqualityErrors(t *testing.T) {
	sentinel := errors.New("equality failed")
	value := failingEqualityValue{err: sentinel}
	input := runtime.NewArrayWith(value)

	for _, tc := range []struct {
		run  func(context.Context) (runtime.Value, error)
		name string
	}{
		{name: "intersection", run: func(ctx context.Context) (runtime.Value, error) {
			return arrays.Intersection(ctx, input, input)
		}},
		{name: "outersection", run: func(ctx context.Context) (runtime.Value, error) {
			return arrays.Outersection(ctx, input, input)
		}},
		{name: "minus", run: func(ctx context.Context) (runtime.Value, error) {
			return arrays.Minus(ctx, input, input)
		}},
		{name: "remove values", run: func(ctx context.Context) (runtime.Value, error) {
			return arrays.RemoveValues(ctx, input, input)
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, err := tc.run(context.Background())
			if !errors.Is(err, sentinel) {
				t.Fatalf("expected equality error, got %v", err)
			}
		})
	}
}

func assertListContainsExactly(t *testing.T, ctx context.Context, list runtime.List, expected ...runtime.Value) {
	t.Helper()

	length, err := list.Length(ctx)
	if err != nil {
		t.Fatalf("length: %v", err)
	}
	if int(length) != len(expected) {
		t.Fatalf("expected %d values, got %d", len(expected), length)
	}

	matched := make([]bool, len(expected))
	err = list.ForEach(ctx, func(ctx context.Context, value runtime.Value, _ runtime.Int) (runtime.Boolean, error) {
		for idx, candidate := range expected {
			if matched[idx] {
				continue
			}
			equal, err := runtime.EqualValues(ctx, value, candidate)
			if err != nil {
				return false, err
			}
			if equal {
				matched[idx] = true
				return true, nil
			}
		}

		return false, errors.New("unexpected result value")
	})
	if err != nil {
		t.Fatalf("iterate: %v", err)
	}
}
