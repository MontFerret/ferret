package vm

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/internal/valueset"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func arrayAll(ctx context.Context, cmp arrayComparator, left, right runtime.Value) (runtime.Boolean, error) {
	arr, err := runtime.CastList(left)

	if err != nil {
		return runtime.False, err
	}

	pred := cmp.predicate()
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

func arrayAny(ctx context.Context, cmp arrayComparator, left, right runtime.Value) (runtime.Boolean, error) {
	arr, err := runtime.CastList(left)

	if err != nil {
		return runtime.False, err
	}

	pred := cmp.predicate()
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

func arrayNone(ctx context.Context, cmp arrayComparator, left, right runtime.Value) (runtime.Boolean, error) {
	arr, err := runtime.CastList(left)

	if err != nil {
		return runtime.False, err
	}

	pred := cmp.predicate()
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

func arrayFlatten(ctx context.Context, value runtime.Value, depth int) (runtime.List, error) {
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

func arrayDistinct(ctx context.Context, value runtime.Value) (runtime.List, error) {
	list, err := runtime.CastList(value)
	if err != nil {
		return nil, err
	}

	size, err := list.Length(ctx)
	if err != nil {
		return nil, err
	}

	result := runtime.NewArray64(size)
	seen := valueset.New(int(size))

	err = list.ForEach(ctx, func(ctx context.Context, item runtime.Value, _ runtime.Int) (runtime.Boolean, error) {
		if seen.Add(item) {
			if err := result.Append(ctx, item); err != nil {
				return runtime.False, err
			}
		}

		return runtime.True, nil
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}
