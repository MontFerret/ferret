package arrays

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/internal"

	"github.com/MontFerret/ferret/pkg/runtime/core"
)

// UNION returns the union of all passed arrays.
// @param {Any[], repeated} arrays - List of arrays to combine.
// @return {Any[]} - All array elements combined in a single array, in any order.
func Union(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 2, core.MaxArgs)

	if err != nil {
		return core.None, err
	}

	err = core.AssertList(args[0])

	if err != nil {
		return core.None, err
	}

	firstArrLen := args[0].(*internal.Array).Length()
	result := internal.NewArray(len(args) * int(firstArrLen))

	for _, arg := range args {
		err := core.AssertList(arg)

		if err != nil {
			return core.None, err
		}

		arr := arg.(*internal.Array)

		arr.ForEach(func(value core.Value, _ int) bool {
			result.Push(value)
			return true
		})
	}

	return result, nil
}
