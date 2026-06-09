package arrays

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/internal/valueset"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// UNION_DISTINCT returns the union of all passed arrays with unique values.
// @param {Any[], repeated} arrays - List of arrays to combine.
// @return {Any[]} - All unique array elements combined in a single array, in any order.
func UnionDistinct(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgs(args, 2, runtime.MaxArgs); err != nil {
		return runtime.None, err
	}

	list, err := runtime.CastArgAt[runtime.List](args, 0)

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

	seen := valueset.New(0)
	result := runtime.NewArray(capacity)

	for i, arg := range args {
		currList, err := runtime.CastArg[runtime.List](arg, i)

		if err != nil {
			return runtime.None, err
		}

		err = currList.ForEach(ctx, func(ctx context.Context, value runtime.Value, idx runtime.Int) (runtime.Boolean, error) {
			if !seen.Add(value) {
				return true, nil
			}

			return true, result.Append(ctx, value)
		})

		if err != nil {
			return runtime.None, err
		}
	}

	return result, nil
}
