package operators

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime"
)

func ArrayAll(ctx context.Context, cmp Comparator, left, right runtime.Value) (runtime.Boolean, error) {
	arr, err := runtime.CastList(left)

	if err != nil {
		return runtime.False, err
	}

	pred := cmp.Predicate()
	result := runtime.True

	if err := arr.ForEach(ctx, func(ctx context.Context, v runtime.Value, _ runtime.Int) (runtime.Boolean, error) {
		if !pred(ctx, v, right) {
			result = runtime.False
			return runtime.False, nil
		}

		return runtime.True, nil
	}); err != nil {
		return runtime.False, err
	}

	return result, nil
}

func ArrayAny(ctx context.Context, cmp Comparator, left, right runtime.Value) (runtime.Boolean, error) {
	arr, err := runtime.CastList(left)

	if err != nil {
		return runtime.False, err
	}

	pred := cmp.Predicate()
	result := runtime.False

	if err := arr.ForEach(ctx, func(ctx context.Context, v runtime.Value, _ runtime.Int) (runtime.Boolean, error) {
		if pred(ctx, v, right) {
			result = runtime.True
			return runtime.False, nil
		}
		return runtime.True, nil
	}); err != nil {
		return runtime.False, err
	}

	return result, nil
}

func ArrayNone(ctx context.Context, cmp Comparator, left, right runtime.Value) (runtime.Boolean, error) {
	arr, err := runtime.CastList(left)

	if err != nil {
		return runtime.False, err
	}

	pred := cmp.Predicate()
	result := runtime.True

	if err := arr.ForEach(ctx, func(ctx context.Context, v runtime.Value, _ runtime.Int) (runtime.Boolean, error) {
		if pred(ctx, v, right) {
			result = runtime.False
			return runtime.False, nil
		}
		return runtime.True, nil
	}); err != nil {
		return runtime.False, err
	}

	return result, nil
}
