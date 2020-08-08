package arrays

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

// PUSH create a new array with appended value.
// @param {Any[]} array - Source array.
// @param {Any} value - Target value.
// @param {Boolean} [unique=False] - Read indicating whether to do uniqueness check.
// @return {Any[]} - A new array with appended value.
func Push(_ context.Context, args ...core.Value) (core.Value, error) {
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
	uniq := false

	if len(args) > 2 {
		err = core.ValidateType(args[2], types.Boolean)

		if err != nil {
			return values.None, err
		}

		uniq = args[2].Compare(values.True) == 0
	}

	result := values.NewArray(int(arr.Length() + 1))
	push := true

	arr.ForEach(func(item core.Value, idx int) bool {
		if uniq && push {
			push = item.Compare(value) != 0
		}

		result.Push(item)

		return true
	})

	if push {
		result.Push(value)
	}

	return result, nil
}
