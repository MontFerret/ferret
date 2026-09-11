package arrays

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// PushMutable appends one value to the original array, retaining duplicates.
// @param array {Any[]} Array to mutate.
// @param value {Any} Value to insert as one element.
// @return {Any[]} The original array.
func PushMutable(ctx context.Context, array, value runtime.Value) (runtime.Value, error) {
	list, err := mutableListArg(ctx, array)
	if err != nil {
		return runtime.None, err
	}

	return mutationResult(array, list.Append(ctx, value))
}
