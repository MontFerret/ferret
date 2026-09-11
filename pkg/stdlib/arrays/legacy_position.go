package arrays

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// Reports whether an element is contained in the array.
// @param array {Any[]} The source array.
// @param value {Any} The target value.
// @return {Boolean} Whether an equal value is present.
// @deprecated Use arrays::contains or arrays::index_of.
func legacyPosition2(ctx context.Context, arg1, arg2 runtime.Value) (runtime.Value, error) {
	return Contains(ctx, arg1, arg2)
}

// position returns a value indicating whether an element is contained in array. Optionally returns its position.
// @param array {Any[]} The source array.
// @param value {Any} The target value.
// @param position {Boolean} Boolean value which indicates whether to return item's position.
// @return {Boolean | Int} A value indicating whether an element is contained in array.
// @deprecated Use arrays::contains or arrays::index_of.
func legacyPosition3(ctx context.Context, arg1, arg2, arg3 runtime.Value) (runtime.Value, error) {
	arr, err := runtime.CastArg[runtime.List](arg1, 0)
	if err != nil {
		return runtime.None, err
	}

	retIdx, err := runtime.CastArg[runtime.Boolean](arg3, 2)
	if err != nil {
		return runtime.None, err
	}

	if retIdx {
		return IndexOf(ctx, arr, arg2)
	}

	return Contains(ctx, arr, arg2)
}
