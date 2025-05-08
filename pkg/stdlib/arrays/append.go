package arrays

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime"
)

// APPEND appends a new item to an array and returns a new array with a given element.
// If “uniqueOnly“ is set to true, then will add the item only if it's unique.
// @param {Any[]} arr - Target array.
// @param {Any} item - Target value to add.
// @param {Boolean} [unique=false] - If set to true, will add the item only if it's unique.
// @return {Any[]} - New array.
func Append(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgs(args, 2, 3); err != nil {
		return runtime.None, err
	}

	list, err := runtime.CastList(args[0])

	if err != nil {
		return runtime.None, err
	}

	arg := args[1]
	unique := runtime.False

	if len(args) > 2 {
		arg3, err := runtime.CastBoolean(args[2])

		if err != nil {
			return runtime.None, err
		}

		unique = arg3
	}

	// We do not know for sure if the list is an array or custom List implementation.
	// Hence, we must solely rely on the List interface.
	next, err := list.CopyWithCap(ctx, 1)

	if err != nil {
		return runtime.None, err
	}

	if unique {
		contains, err := list.Contains(ctx, arg)

		if err != nil {
			return runtime.None, err
		}

		if contains {
			return next, nil
		}
	}

	if err := next.Add(ctx, arg); err != nil {
		return runtime.None, err
	}

	return next, nil
}
