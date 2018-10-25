package arrays

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

// RemoveValues returns a new array with removed all occurrences of values in a given array.
// @param array (Array) - Source array.
// @param values (Array) - Target values.
// @returns (Array) - A new array with removed all occurrences of values in a given array.
func RemoveValues(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 2, 2)

	if err != nil {
		return values.None, err
	}

	err = core.ValidateType(args[0], types.Array)

	if err != nil {
		return values.None, err
	}

	err = core.ValidateType(args[1], types.Array)

	if err != nil {
		return values.None, err
	}

	arr := args[0].(*values.Array)
	vals := args[1].(*values.Array)

	result := values.NewArray(int(arr.Length()))
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
