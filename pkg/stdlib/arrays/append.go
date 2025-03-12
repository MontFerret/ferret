package arrays

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/internal"

	"github.com/MontFerret/ferret/pkg/runtime/core"
)

// APPEND appends a new item to an array and returns a new array with a given element.
// If “uniqueOnly“ is set to true, then will add the item only if it's unique.
// @param {Any[]} arr - Target array.
// @param {Any} item - Target value to add.
// @param {Boolean} [unique=false] - If set to true, will add the item only if it's unique.
// @return {Any[]} - New array.
func Append(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 2, 3)

	if err != nil {
		return core.None, err
	}

	err = core.AssertList(args[0])

	if err != nil {
		return core.None, err
	}

	arr := args[0].(*internal.Array)
	arg := args[1]
	unique := core.False

	if len(args) > 2 {
		err = core.AssertBoolean(args[2])

		if err != nil {
			return core.None, err
		}

		unique = args[2].(core.Boolean)
	}

	next := internal.NewArray(int(arr.Length()) + 1)

	if !unique {
		arr.ForEach(func(item core.Value, idx int) bool {
			next.Push(item)

			return true
		})

		next.Push(arg)

		return next, nil
	}

	hasDuplicate := false

	arr.ForEach(func(item core.Value, idx int) bool {
		next.Push(item)

		if !hasDuplicate {
			hasDuplicate = core.CompareValues(item, arg) == 0
		}

		return true
	})

	if !hasDuplicate {
		next.Push(arg)
	}

	return next, nil
}
