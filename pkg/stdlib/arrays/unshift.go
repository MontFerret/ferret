package arrays

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

// UNSHIFT prepends value to a given array.
// @param {Any[]} array - Target array.
// @param {Any} value - Target value to prepend.
// @param {Boolean} [unique=False] - Optional value indicating whether a value must be unique to be prepended. Default is false.
// @return {Any[]} - New array with prepended value.
func Unshift(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 2, 3)

	if err != nil {
		return values.None, err
	}

	err = core.ValidateType(args[0], types.Array)

	if err != nil {
		return values.None, err
	}

	arr := args[0].(*values.Array)
	value := args[1]
	uniq := values.False

	if len(args) > 2 {
		err = core.ValidateType(args[2], types.Boolean)

		if err != nil {
			return values.None, err
		}

		uniq = args[2].(values.Boolean)
	}

	result := values.NewArray(int(arr.Length() + 1))

	if !uniq {
		result.Push(value)

		arr.ForEach(func(el core.Value, _ int) bool {
			result.Push(el)

			return true
		})
	} else {
		ok := true

		// let's just hope it's unique
		// if not, we will terminate the loop and return a copy of an array
		result.Push(value)

		arr.ForEach(func(el core.Value, idx int) bool {
			if el.Compare(value) != 0 {
				result.Push(el)

				return true
			}

			// not unique
			ok = false
			return false
		})

		if !ok {
			// value is not unique, just return a new copy with same elements
			return arr.Copy(), nil
		}
	}

	return result, nil
}
