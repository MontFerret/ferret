package objects

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime"
)

// KEEP_KEYS returns a new object with only given keys.
// @param {hashMap} obj - Source object.
// @param {String, repeated} keys - Keys that need to be kept.
// @return {hashMap} - New hashMap with only given keys.
// TODO: REWRITE TO USE LIST & MAP instead
func KeepKeys(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgs(args, 2, runtime.MaxArgs); err != nil {
		return runtime.None, err
	}

	if err := runtime.AssertMap(args[0]); err != nil {
		return runtime.None, err
	}

	var keys *runtime.Array

	if len(args) == 2 {
		arr, ok := args[1].(*runtime.Array)

		if ok {
			keys = arr
		}
	}

	if keys == nil {
		keys = runtime.NewArrayWith(args[1:]...)
	}

	if err := validateArrayOf(ctx, runtime.TypeString, keys); err != nil {
		return runtime.None, err
	}

	obj := args[0].(*runtime.Object)
	resultObj := runtime.NewObject()

	var key runtime.String

	_ = keys.ForEach(ctx, func(c context.Context, keyVal runtime.Value, idx runtime.Int) (runtime.Boolean, error) {
		key = keyVal.(runtime.String)

		if val, err := obj.Get(c, key); err == nil {
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

	return resultObj, nil
}
