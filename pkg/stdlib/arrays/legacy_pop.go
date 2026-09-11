package arrays

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// pop returns a new array without last element.
// @param array {Any[]} Target array.
// @return {Any[]} Copy of an array without last element.
// @deprecated Use arrays::slice with a nonnegative length.
func legacyPop(ctx context.Context, arg runtime.Value) (runtime.Value, error) {
	arr, err := runtime.CastArg[runtime.List](arg, 0)
	if err != nil {
		return runtime.None, err
	}

	size, err := arr.Length(ctx)
	if err != nil {
		return runtime.None, err
	}

	if size == 0 {
		return arr.Empty(ctx)
	}

	return arr.Slice(ctx, 0, size-1)
}
