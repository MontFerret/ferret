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

func push2(ctx context.Context, arg1, arg2 runtime.Value) (runtime.Value, error) {
	return append2(ctx, arg1, arg2)
}

func push3(ctx context.Context, arg1, arg2, arg3 runtime.Value) (runtime.Value, error) {
	return append3(ctx, arg1, arg2, arg3)
}
