package arrays

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

// UNSHIFT prepends value to a given array.
// @param {Any[]} array - Target array.
// @param {Any} value - Target value to prepend.
// @param {Boolean} [unique=False] - Optional value indicating whether a value must be unique to be prepended. Default is false.
// @return {Any[]} - New array with prepended value.
func Unshift(_ context.Context, args ...core.Value) (core.Value, error) {
	if err := core.ValidateArgs(args, 2, 3); err != nil {
		return values.None, err
	}

	arr, err := values.CastArray(args[0])
	value := args[1]
	uniq := values.False

	if len(args) > 2 {
		uniq, err = values.CastBoolean(args[2])

		if err != nil {
			return values.None, err
		}
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
			if values.Compare(el, value) != 0 {
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
