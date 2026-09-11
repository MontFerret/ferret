package arrays

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// RemoveAtMutable removes and returns an existing element.
// Negative and missing indexes fail before mutation.
// @param array {Any[]} Array to mutate.
// @param index {Int} Zero-based element index.
// @return {Any} Removed value.
func RemoveAtMutable(ctx context.Context, array, index runtime.Value) (runtime.Value, error) {
	list, position, err := mutableIndexArgs(ctx, array, index, false)
	if err != nil {
		return runtime.None, err
	}

	return removeMutableElement(ctx, list, position)
}
