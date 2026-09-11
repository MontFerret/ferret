package arrays

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// ShiftMutable removes and returns the first element, or none when empty.
// @param array {Any[]} Array to mutate.
// @return {Any} Removed value, or none for an empty array.
func ShiftMutable(ctx context.Context, array runtime.Value) (runtime.Value, error) {
	return removeEndMutable(ctx, array, false)
}
