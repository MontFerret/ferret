package arrays

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime"
)

// SORTED_UNIQUE sorts all elements in a array.
// The function will use the default comparison order for FQL value types.
// Additionally, the values in the result array will be made unique
// @param {Any[]} array - Target array.
// @return {Any[]} - Sorted array.
func SortedUnique(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgs(args, 1, 1); err != nil {
		return runtime.None, err
	}

	list, err := runtime.CastList(args[0])

	if err != nil {
		return runtime.None, err
	}

	uniq, err := ToUniqueList(ctx, list)

	if err != nil {
		return runtime.None, err
	}

	if err := uniq.SortAsc(ctx); err != nil {
		return runtime.None, err
	}

	return uniq, nil
}
