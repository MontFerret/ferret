package arrays

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// PUSH create a new array with appended value.
// @param {Any[]} array - Source array.
// @param {Any} value - Target value.
// @param {Boolean} [unique=False] - Read indicating whether to do uniqueness check.
// @return {Any[]} - A new array with appended value.
func Push(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
	return Append(ctx, args...)
}
