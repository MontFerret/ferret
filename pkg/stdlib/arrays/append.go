package arrays

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

// APPEND appends a new item to an array and returns a new array with a given element.
// If ``uniqueOnly`` is set to true, then will add the item only if it's unique.
// @param {Any[]} arr - Target array.
// @param {Any} item - Target value to add.
// @return {Any[]} - New array.
func Append(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 2, 3)

	if err != nil {
		return values.None, err
	}

	err = core.ValidateType(args[0], types.Array)

	if err != nil {
		return values.None, err
	}

	arr := args[0].(*values.Array)
	arg := args[1]
	unique := values.False

	if len(args) > 2 {
		err = core.ValidateType(args[2], types.Boolean)

		if err != nil {
			return values.None, err
		}

		unique = args[2].(values.Boolean)
	}

	next := values.NewArray(int(arr.Length()) + 1)

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
			hasDuplicate = item.Compare(arg) == 0
		}

		return true
	})

	if !hasDuplicate {
		next.Push(arg)
	}

	return next, nil
}
