package objects

import (
	"context"
	"fmt"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// Zip constructs an object from parallel key and value lists of equal length.
// Keys must be strings. Duplicate keys use the last value in the input.
// @param keys {String[]} List of object keys.
// @param values {Any[]} Parallel list of values.
// @return {Map} Independent object containing cloned or copied values.
func Zip(ctx context.Context, arg1, arg2 runtime.Value) (runtime.Value, error) {
	keys, vals, err := runtime.CastArgs2[runtime.List, runtime.List](arg1, arg2)
	if err != nil {
		return runtime.None, err
	}

	if err := ctx.Err(); err != nil {
		return runtime.None, err
	}

	keysSize, err := keys.Length(ctx)
	if err != nil {
		return runtime.None, runtime.ArgError(err, 0)
	}

	valsSize, err := vals.Length(ctx)
	if err != nil {
		return runtime.None, runtime.ArgError(err, 1)
	}

	if keysSize < 0 {
		return runtime.None, runtime.ArgError(runtime.Error(runtime.ErrInvalidArgument, "list length must not be negative"), 0)
	}

	if valsSize < 0 {
		return runtime.None, runtime.ArgError(runtime.Error(runtime.ErrInvalidArgument, "list length must not be negative"), 1)
	}

	if keysSize != valsSize {
		return runtime.None, runtime.ArgError(runtime.Error(runtime.ErrInvalidArgument,
			fmt.Sprintf("keys and values must have the same length; got keys: %d, values: %d", keysSize, valsSize)), 1)
	}

	result := runtime.NewObject()
	for index := runtime.ZeroInt; index < keysSize; index++ {
		if err := ctx.Err(); err != nil {
			return runtime.None, err
		}

		keyValue, err := keys.At(ctx, index)
		if err != nil {
			return runtime.None, runtime.ArgError(fmt.Errorf("item %d: %w", index, err), 0)
		}

		key, err := runtime.CastString(keyValue)
		if err != nil {
			return runtime.None, runtime.ArgError(fmt.Errorf("item %d: %w", index, err), 0)
		}

		value, err := vals.At(ctx, index)
		if err != nil {
			return runtime.None, runtime.ArgError(fmt.Errorf("item %d: %w", index, err), 1)
		}

		if err := setCopiedEntry(ctx, result, key, value); err != nil {
			return runtime.None, runtime.ArgError(fmt.Errorf("item %d: %w", index, err), 1)
		}
	}

	return result, nil
}
