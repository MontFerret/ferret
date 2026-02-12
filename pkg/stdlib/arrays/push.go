package arrays

import (
	"github.com/MontFerret/ferret/pkg/runtime"
)

// PUSH create a new array with appended value.
// @param {Any[]} array - Source array.
// @param {Any} value - Target value.
// @param {Boolean} [unique=False] - Read indicating whether to do uniqueness check.
// @return {Any[]} - A new array with appended value.
func Push(ctx runtime.Context, args ...runtime.Value) (runtime.Value, error) {
	// TODO: Why do we have two functions with the same functionality?
	return Append(ctx, args...)
}
