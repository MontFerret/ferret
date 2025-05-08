package arrays

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime"
)

// OUTERSECTION return the values that occur only once across all arrays specified.
// The element order is random.
// @param {Any[], repeated} arrays - An arbitrary number of arrays as multiple arguments (at least 2).
// @return {Any[]} - A single array with only the elements that exist only once across all provided arrays.
func Outersection(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
	return sections(ctx, args, 1)
}
