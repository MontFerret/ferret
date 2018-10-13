package arrays

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
)

// Outersection return the values that occur only once across all arrays specified.
// @param arrays (Array, repeated) - An arbitrary number of arrays as multiple arguments (at least 2).
// @returns (Array) - A single array with only the elements that exist only once across all provided arrays.
// The element order is random.
func Outersection(_ context.Context, args ...core.Value) (core.Value, error) {
	return sections(args, 1)
}
