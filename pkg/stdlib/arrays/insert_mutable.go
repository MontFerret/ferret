package arrays

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// InsertMutable inserts an element, allowing index equal to length.
// Invalid indexes fail before mutation.
// @param array {Any[]} Array to mutate.
// @param index {Int} Zero-based insertion index.
// @param value {Any} Value to insert as one element.
// @return {Any[]} The original array.
func InsertMutable(ctx context.Context, array, index, value runtime.Value) (runtime.Value, error) {
	list, position, err := mutableIndexArgs(ctx, array, index, true)
	if err != nil {
		return runtime.None, err
	}

	return mutationResult(array, list.Insert(ctx, position, value))
}
