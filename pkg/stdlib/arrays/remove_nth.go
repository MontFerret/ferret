package arrays

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/internal"

	"github.com/MontFerret/ferret/pkg/runtime/core"
)

// REMOVE_NTH returns a new array without an element by a given position.
// @param {Any[]} array - Source array.
// @param {Int} position - Target element position.
// @return {Any[]} - A new array without an element by a given position.
func RemoveNth(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 2, 2)

	if err != nil {
		return core.None, err
	}

	err = core.AssertList(args[0])

	if err != nil {
		return core.None, err
	}

	err = core.AssertInt(args[1])

	if err != nil {
		return core.None, err
	}

	arr := args[0].(*internal.Array)
	index := int(args[1].(core.Int))
	result := internal.NewArray(int(arr.Length() - 1))

	arr.ForEach(func(value core.Value, idx int) bool {
		if idx != index {
			result.Push(value)
		}

		return true
	})

	return result, nil
}
