package arrays

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// remove_values returns a new array with removed all occurrences of values in a given array.
// @param array {Any[]} Source array.
// @param values {Any[]} Target values.
// @return {Any[]} A new array with removed all occurrences of values in a given array.
func RemoveValues(ctx context.Context, arg1, arg2 runtime.Value) (runtime.Value, error) {
	arr, vals, err := runtime.CastArgs2[runtime.List, runtime.List](arg1, arg2)

	if err != nil {
		return runtime.None, err
	}

	lookupTable := make(map[uint64][]runtime.Value)

	err = vals.ForEach(ctx, func(ctx context.Context, value runtime.Value, idx runtime.Int) (runtime.Boolean, error) {
		hash := value.Hash()
		bucket := lookupTable[hash]
		match, err := findEqualValue(ctx, bucket, value)
		if err != nil {
			return false, err
		}

		if match < 0 {
			lookupTable[hash] = append(bucket, value)
		}

		return true, nil
	})

	if err != nil {
		return runtime.None, err
	}

	return arr.Filter(ctx, func(ctx context.Context, value runtime.Value, idx runtime.Int) (runtime.Boolean, error) {
		h := value.Hash()

		match, err := findEqualValue(ctx, lookupTable[h], value)
		if err != nil {
			return false, err
		}

		return match < 0, nil
	})
}
