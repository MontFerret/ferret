package arrays

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

/*
 * Returns the union of all passed arrays.
 * @param arrays (Array, repeated) - List of arrays to combine.
 * @returns (Array) - All array elements combined in a single array, in any order.
 */
func Union(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 2, core.MaxArgs)

	if err != nil {
		return values.None, err
	}

	result := values.NewArray(len(args) * 2)

	for _, arg := range args {
		err := core.ValidateType(arg, core.ArrayType)

		if err != nil {
			return values.None, err
		}

		arr := arg.(*values.Array)

		arr.ForEach(func(value core.Value, _ int) bool {
			result.Push(value)
			return true
		})
	}

	return result, nil
}
