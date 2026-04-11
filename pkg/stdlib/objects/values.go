package objects

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// VALUES return the attribute values of the map as a list.
// @param {Map} map - Target map.
// @return {Any[]} - Values of document returned in any order.
func Values(ctx context.Context, arg runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgType(arg, 0, runtime.TypeMap); err != nil {
		return runtime.None, err
	}

	obj := arg.(runtime.Map)

	values, err := obj.Values(ctx)
	if err != nil {
		return runtime.None, err
	}

	length, err := values.Length(ctx)
	if err != nil {
		return runtime.None, err
	}

	result := runtime.NewArray64(length)

	if err := values.ForEach(ctx, func(c context.Context, value runtime.Value, _ runtime.Int) (runtime.Boolean, error) {
		cloned, err := runtime.CloneOrCopy(c, value)
		if err != nil {
			return false, err
		}

		if err := result.Append(c, cloned); err != nil {
			return false, err
		}

		return true, nil
	}); err != nil {
		return runtime.None, err
	}

	return result, nil
}
