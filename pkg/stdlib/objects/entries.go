package objects

import (
	"context"
	"fmt"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// Entries returns associated key/value pairs from one map traversal.
// Keys and values follow their clone or copy contracts; ordering is unspecified.
// @param value {Map} Map whose entries are returned.
// @return {Any[][]} Independent two-item lists containing each key and its value.
func Entries(ctx context.Context, arg runtime.Value) (runtime.Value, error) {
	src, err := runtime.CastArg[runtime.Map](arg, 0)
	if err != nil {
		return runtime.None, err
	}

	if err := ctx.Err(); err != nil {
		return runtime.None, err
	}

	result := runtime.NewArray(0)
	err = src.ForEach(ctx, func(ctx context.Context, value, key runtime.Value) (runtime.Boolean, error) {
		if err := ctx.Err(); err != nil {
			return false, err
		}

		copiedKey, err := runtime.CloneOrCopy(ctx, key)
		if err != nil {
			return false, fmt.Errorf("key %q: %w", key.String(), err)
		}

		copiedValue, err := runtime.CloneOrCopy(ctx, value)
		if err != nil {
			return false, fmt.Errorf("key %q: %w", key.String(), err)
		}

		if err := result.Append(ctx, runtime.NewArrayWith(copiedKey, copiedValue)); err != nil {
			return false, err
		}

		return true, nil
	})
	if err != nil {
		return runtime.None, runtime.ArgError(err, 0)
	}

	return result, nil
}
