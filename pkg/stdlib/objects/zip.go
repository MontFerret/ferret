package objects

import (
	"context"
	"fmt"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// ZIP returns an object assembled from the separate parameters keys and values.
// Keys and values must be arrays and have the same length.
// @param {String[]} keys - An array of strings, to be used as key names in the result.
// @param {hashMap[]} values - An array of runtime.Value, to be used as key values.
// @return {Map} - An object with the keys and values assembled.
func Zip(ctx context.Context, arg1, arg2 runtime.Value) (runtime.Value, error) {
	keys, err := runtime.CastArg[runtime.List](arg1, 0)

	if err != nil {
		return runtime.None, err
	}

	vals, err := runtime.CastArg[runtime.List](arg2, 1)

	if err != nil {
		return runtime.None, err
	}

	keysSize, _ := keys.Length(ctx)
	valsSize, _ := vals.Length(ctx)

	if keysSize != valsSize {
		return runtime.None, runtime.Error(
			runtime.ErrInvalidArgument,
			fmt.Sprintf("keys and values must have the same length. got keys: %d, values: %d",
				keysSize, valsSize,
			),
		)
	}

	if err := runtime.AssertItemsOf(ctx, keys, runtime.AssertString); err != nil {
		return runtime.None, runtime.ArgError(err, 1)
	}

	zipped := runtime.NewObject()

	var k runtime.String
	var val runtime.Value
	var exists bool
	keyExists := map[runtime.String]bool{}

	return zipped, keys.ForEach(ctx, func(c context.Context, key runtime.Value, idx runtime.Int) (runtime.Boolean, error) {
		k = key.(runtime.String)

		// If the key already exists, we skip it. This allows us to handle duplicate keys in the input.
		if _, exists = keyExists[k]; exists {
			return true, nil
		}

		keyExists[k] = true

		val, _ = vals.At(c, idx)

		cloneable, ok := val.(runtime.Cloneable)

		if ok {
			val, _ = cloneable.Clone(c)
		}

		_ = zipped.Set(c, k, val)

		return true, nil
	})
}
