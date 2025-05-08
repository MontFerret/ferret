package arrays

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime"
)

// REMOVE_VALUES returns a new array with removed all occurrences of values in a given array.
// @param {Any[]} array - Source array.
// @param {Any[]} values - Target values.
// @return {Any[]} - A new array with removed all occurrences of values in a given array.
func RemoveValues(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgs(args, 2, 2); err != nil {
		return runtime.None, err
	}

	arr, err := runtime.CastList(args[0])

	if err != nil {
		return runtime.None, err
	}

	vals, err := runtime.CastList(args[1])

	if err != nil {
		return runtime.None, err
	}

	lookupTable := make(map[uint64]bool)

	err = vals.ForEach(ctx, func(ctx context.Context, value runtime.Value, idx runtime.Int) (runtime.Boolean, error) {
		lookupTable[value.Hash()] = true

		return true, nil
	})

	if err != nil {
		return runtime.None, err
	}

	return arr.Find(ctx, func(ctx context.Context, value runtime.Value, idx runtime.Int) (runtime.Boolean, error) {
		h := value.Hash()

		_, exists := lookupTable[h]

		return runtime.Boolean(!exists), nil
	})
}
