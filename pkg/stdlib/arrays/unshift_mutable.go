package arrays

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// UnshiftMutable inserts one value at the beginning of the original array.
// @param array {Any[]} Array to mutate.
// @param value {Any} Value to insert as one element.
// @return {Any[]} The original array.
func UnshiftMutable(ctx context.Context, array, value runtime.Value) (runtime.Value, error) {
	return InsertMutable(ctx, array, runtime.ZeroInt, value)
}
