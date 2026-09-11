package arrays

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// SortMutable stably sorts the original array in ascending FQL order.
// It uses the same ordering as arrays::sorted; errors may leave a partial sort.
// @param array {Any[]} Array to mutate.
// @return {Any[]} The original array.
func SortMutable(ctx context.Context, array runtime.Value) (runtime.Value, error) {
	list, err := mutableListArg(ctx, array)
	if err != nil {
		return runtime.None, err
	}

	if _, err := mutableListLength(ctx, list); err != nil {
		return runtime.None, runtime.ArgError(err, 0)
	}

	return mutationResult(array, runtime.SortAsc(ctx, list))
}
