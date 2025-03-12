package arrays

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/internal"

	"github.com/MontFerret/ferret/pkg/runtime/core"
)

// REMOVE_VALUES returns a new array with removed all occurrences of values in a given array.
// @param {Any[]} array - Source array.
// @param {Any[]} values - Target values.
// @return {Any[]} - A new array with removed all occurrences of values in a given array.
func RemoveValues(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 2, 2)

	if err != nil {
		return core.None, err
	}

	err = core.AssertList(args[0])

	if err != nil {
		return core.None, err
	}

	err = core.AssertList(args[1])

	if err != nil {
		return core.None, err
	}

	arr := args[0].(*internal.Array)
	vals := args[1].(*internal.Array)

	result := internal.NewArray(int(arr.Length()))
	lookupTable := make(map[uint64]bool)

	vals.ForEach(func(value core.Value, idx int) bool {
		lookupTable[value.Hash()] = true

		return true
	})

	arr.ForEach(func(value core.Value, idx int) bool {
		h := value.Hash()

		_, exists := lookupTable[h]

		if !exists {
			result.Push(value)
		}

		return true
	})

	return result, nil
}
