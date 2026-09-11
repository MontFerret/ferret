package arrays

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// RemoveMutable removes all equal values from the original array, preserving survivor order.
// Errors may leave earlier writes or removals applied.
// @param array {Any[]} Array to mutate.
// @param value {Any} Value whose occurrences are removed.
// @return {Any[]} The original array.
func RemoveMutable(ctx context.Context, array, value runtime.Value) (runtime.Value, error) {
	list, err := mutableListArg(ctx, array)
	if err != nil {
		return runtime.None, err
	}

	return mutationResult(array, removeMatchingInPlace(ctx, list, value))
}
