package arrays

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// ClearMutable removes all elements from the original array without closing them.
// @param array {Any[]} Array to mutate.
// @return {Any[]} The original array.
func ClearMutable(ctx context.Context, array runtime.Value) (runtime.Value, error) {
	list, err := mutableListArg(ctx, array)
	if err != nil {
		return runtime.None, err
	}

	return mutationResult(array, list.Clear(ctx))
}
