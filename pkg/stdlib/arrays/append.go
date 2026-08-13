package arrays

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// append appends a new item to an array and returns a new array with a given element.
// If “uniqueOnly“ is set to true, then will add the item only if it's unique.
// @param arr {Any[]} Target array.
// @param item {Any} Target value to add.
// @param unique {Boolean} If set to true, will add the item only if it's unique.
// @return {Any[]} New array.
func Append(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgs(args, 2, 3); err != nil {
		return runtime.None, err
	}

	if len(args) == 2 {
		return append2(ctx, args[0], args[1])
	}

	return append3(ctx, args[0], args[1], args[2])
}

// append appends a new item to an array and returns a new array with a given element.
// If “uniqueOnly“ is set to true, then will add the item only if it's unique.
// @param arr {Any[]} Target array.
// @param item {Any} Target value to add.
// @return {Any[]} New array.
func append2(ctx context.Context, arg1, arg2 runtime.Value) (runtime.Value, error) {
	return append3(ctx, arg1, arg2, runtime.False)
}

// append appends a new item to an array and returns a new array with a given element.
// If “uniqueOnly“ is set to true, then will add the item only if it's unique.
// @param arr {Any[]} Target array.
// @param item {Any} Target value to add.
// @param unique {Boolean} If set to true, will add the item only if it's unique.
// @return {Any[]} New array.
func append3(ctx context.Context, arg1, arg2, arg3 runtime.Value) (runtime.Value, error) {
	list, err := runtime.CastArg[runtime.List](arg1, 0)

	if err != nil {
		return runtime.None, err
	}

	unique, err := runtime.CastArg[runtime.Boolean](arg3, 2)
	if err != nil {
		return runtime.None, err
	}

	var next runtime.List

	// We do not know for sure if the list is an array or custom List implementation.
	// Hence, we must solely rely on the List interface.
	switch v := list.(type) {
	case *runtime.Array:
		next = v.CopyWithGrowth(1)
	case runtime.List:
		next = v.Copy().(runtime.List)
	}

	if unique {
		idx, err := list.IndexOf(ctx, arg2)

		if err != nil {
			return runtime.None, err
		}

		if idx > -1 {
			return next, nil
		}
	}

	if err := next.Append(ctx, arg2); err != nil {
		return runtime.None, err
	}

	return next, nil
}
