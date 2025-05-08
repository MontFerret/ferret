package arrays

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime"
)

// UNION returns the union of all passed arrays.
// @param {Any[], repeated} arrays - List of arrays to combine.
// @return {Any[]} - All array elements combined in a single array, in any order.
func Union(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgs(args, 2, runtime.MaxArgs); err != nil {
		return runtime.None, err
	}

	list, err := runtime.CastList(args[0])

	if err != nil {
		return runtime.None, err
	}

	firstSize, err := list.Length(ctx)

	if err != nil {
		return runtime.None, err
	}

	capacity := len(args) * int(firstSize)

	if capacity == 0 {
		capacity = len(args) * 5
	}

	result := runtime.NewArray(capacity)

	for _, arg := range args {
		currList, err := runtime.CastList(arg)

		if err != nil {
			return runtime.None, err
		}

		err = currList.ForEach(ctx, func(ctx context.Context, value runtime.Value, idx runtime.Int) (runtime.Boolean, error) {
			return true, result.Add(ctx, value)
		})

		if err != nil {
			return runtime.None, err
		}
	}

	return result, nil
}
