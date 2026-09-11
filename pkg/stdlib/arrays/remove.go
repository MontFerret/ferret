package arrays

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// Remove returns a shallow copy without any values equal to the target.
// @param array {Any[]} Source array.
// @param value {Any} Value whose occurrences are removed.
// @return {Any[]} Remaining elements in their original order.
func Remove(ctx context.Context, array, value runtime.Value) (runtime.Value, error) {
	list, err := runtime.CastArg[runtime.List](array, 0)
	if err != nil {
		return runtime.None, err
	}

	return removeLimited(ctx, list, value, -1)
}
