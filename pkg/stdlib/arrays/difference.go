package arrays

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// Difference returns distinct values from the first array absent from every later array.
// @param arrays {Any[], repeated} At least two arrays.
// @return {Any[]} Remaining values in first-input order.
func Difference(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgs(args, 2, runtime.MaxArgs); err != nil {
		return runtime.None, err
	}

	first, err := runtime.CastArgAt[runtime.List](args, 0)
	if err != nil {
		return runtime.None, err
	}

	candidates, ordered, err := collectSet(ctx, first)
	if err != nil {
		return runtime.None, err
	}

	for index, arg := range args[1:] {
		list, err := runtime.CastArg[runtime.List](arg, index+1)
		if err != nil {
			return runtime.None, err
		}

		err = list.ForEach(ctx, func(ctx context.Context, value runtime.Value, _ runtime.Int) (runtime.Boolean, error) {
			_, err := candidates.Remove(ctx, value)

			return true, err
		})
		if err != nil {
			return runtime.None, err
		}
	}

	return filterSet(ctx, ordered, candidates, true)
}
