package objects

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// KEEP_KEYS returns a new object with only given keys.
// @param {Object|Map} obj - Source object.
// @param {String, repeated} keys - Keys that need to be kept.
// @return {Map} - New Object with only given keys.
func KeepKeys(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgs(args, 2, runtime.MaxArgs); err != nil {
		return runtime.None, err
	}

	if err := runtime.ValidateArgTypeAt(args, 0, runtime.TypeMap); err != nil {
		return runtime.None, err
	}

	src := args[0].(runtime.Map)

	var keys runtime.List

	if len(args) == 2 {
		list, ok := args[1].(runtime.List)

		if ok {
			keys = list
		}
	}

	if keys == nil {
		keys = runtime.NewArrayWith(args[1:]...)
	}

	if err := runtime.AssertItemsOf(ctx, keys, runtime.AssertString); err != nil {
		return runtime.None, runtime.ArgError(err, 1)
	}

	resultObj, err := src.Empty(ctx)

	if err != nil {
		return runtime.None, err
	}

	var key runtime.String

	return resultObj, keys.ForEach(ctx, func(c context.Context, keyVal runtime.Value, idx runtime.Int) (runtime.Boolean, error) {
		key = keyVal.(runtime.String)

		if val, err := src.Get(c, key); err == nil {
			cloneable, ok := val.(runtime.Cloneable)

			if ok {
				v, err := cloneable.Clone(c)

				if err != nil {
					return runtime.False, err
				}

				val = v
			}

			_ = resultObj.Set(c, key, val)
		}

		return true, nil
	})
}
