package arrays

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime"
)

// SORTED sorts all elements in the given list.
// The function will use the default comparison order for FQL value types.
// @param {Any[]} array - Target array.
// @return {Any[]} - Sorted array.
func Sorted(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgs(args, 1, 1); err != nil {
		return runtime.None, err
	}

	list, err := runtime.CastList(args[0])

	if err != nil {
		return runtime.None, err
	}

	size, err := list.Length(ctx)

	if err != nil {
		return runtime.None, err
	}

	if size == 0 {
		return runtime.EmptyArray(), nil
	}

	copied := list.Copy()

	//if err != nil {
	//	return runtime.None, err
	//}

	if err := runtime.SortAsc(ctx, copied); err != nil {
		return runtime.None, err
	}

	return copied, nil
}
