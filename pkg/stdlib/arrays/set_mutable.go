package arrays

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// SetMutable replaces an existing element without growing the array.
// Invalid indexes fail before mutation.
// @param array {Any[]} Array to mutate.
// @param index {Int} Zero-based element index.
// @param value {Any} Value to insert as one element.
// @return {Any[]} The original array.
func SetMutable(ctx context.Context, array, index, value runtime.Value) (runtime.Value, error) {
	list, position, err := mutableIndexArgs(ctx, array, index, false)
	if err != nil {
		return runtime.None, err
	}

	return mutationResult(array, list.SetAt(ctx, position, value))
}
