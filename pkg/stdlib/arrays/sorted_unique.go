package arrays

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/collections"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

// SortedUnique sorts all elements in anyArray.
// The function will use the default comparison order for FQL value types.
// Additionally, the values in the result array will be made unique
// @param array (Array) - Target array.
// @returns (Array) - Sorted array.
func SortedUnique(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, 1)

	if err != nil {
		return values.None, err
	}

	err = core.ValidateType(args[0], core.ArrayType)

	if err != nil {
		return values.None, err
	}

	arr := args[0].(*values.Array)

	if arr.Length() == 0 {
		return values.NewArray(0), nil
	}

	sorter, err := collections.NewSorter(func(first collections.DataSet, second collections.DataSet) (int, error) {
		return first.Get(iteratorValVar).Compare(second.Get(iteratorValVar)), nil
	}, collections.SortDirectionAsc)

	if err != nil {
		return values.None, err
	}

	uniqIterator, err := collections.NewUniqueIterator(
		toArrayIterator(arr),
		iteratorValVar,
	)

	if err != nil {
		return values.None, err
	}

	iterator, err := collections.NewSortIterator(
		uniqIterator,
		sorter,
	)

	if err != nil {
		return values.None, err
	}

	return toArray(iterator)
}
