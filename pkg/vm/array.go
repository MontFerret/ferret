package vm

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
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

func Flatten(ctx context.Context, value runtime.Value, depth int) (runtime.List, error) {
	list, err := runtime.CastList(value)
	if err != nil {
		return nil, err
	}

	if depth < 1 {
		return list, nil
	}

	size, err := list.Length(ctx)
	if err != nil {
		return nil, err
	}

	result := runtime.NewArray64(size * 2)

	var flatten func(input runtime.List, level int) error
	flatten = func(input runtime.List, level int) error {
		return input.ForEach(ctx, func(c context.Context, v runtime.Value, _ runtime.Int) (runtime.Boolean, error) {
			if listValue, ok := v.(runtime.List); ok && level <= depth {
				if err := flatten(listValue, level+1); err != nil {
					return false, err
				}

				return true, nil
			}

			_ = result.Append(c, v)

			return true, nil
		})
	}

	if err := flatten(list, 1); err != nil {
		return nil, err
	}

	return result, nil
}
