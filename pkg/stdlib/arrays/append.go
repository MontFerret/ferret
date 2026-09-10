package arrays

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// Append returns a shallow copy with one additional element, retaining duplicates.
// @param array {Any[]} Source array.
// @param value {Any} Value to append as one element.
// @return {Any[]} New array with the appended value.
func Append(ctx context.Context, array, value runtime.Value) (runtime.Value, error) {
	list, err := runtime.CastArg[runtime.List](array, 0)
	if err != nil {
		return runtime.None, err
	}

	next := copyForAppend(list)
	if err := next.Append(ctx, value); err != nil {
		return runtime.None, err
	}

	return next, nil
}

func copyForAppend(list runtime.List) runtime.List {
	if array, ok := list.(*runtime.Array); ok {
		return array.CopyWithGrowth(1)
	}

	return list.Copy().(runtime.List)
}
