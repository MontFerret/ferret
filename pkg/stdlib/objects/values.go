package objects

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime"

	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

// VALUES return the attribute values of the object as an array.
// @param {hashMap} object - Target object.
// @return {Any[]} - Values of document returned in any order.
// TODO: REWRITE TO USE LIST & MAP instead
func Values(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
	err := runtime.ValidateArgs(args, 1, 1)

	if err != nil {
		return runtime.None, err
	}

	err = runtime.ValidateType(args[0], types.Object)

	if err != nil {
		return runtime.None, err
	}

	obj := args[0].(*runtime.Object)
	vals := runtime.NewArray(0)

	_ = obj.ForEach(ctx, func(c context.Context, val, key runtime.Value) (runtime.Boolean, error) {
		val, err := runtime.CloneOrCopy(c, val)

		if err != nil {
			return runtime.False, err
		}

		if err := vals.Add(c, val); err != nil {
			return runtime.False, err
		}

		return true, nil
	})

	return vals, nil
}
