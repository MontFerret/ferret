package objects

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

// KEEP_KEYS returns a new object with only given keys.
// @param {hashMap} obj - Source object.
// @param {String, repeated} keys - Keys that need to be kept.
// @return {hashMap} - New hashMap with only given keys.
// TODO: REWRITE TO USE LIST & MAP instead
func KeepKeys(ctx context.Context, args ...core.Value) (core.Value, error) {
	if err := core.ValidateArgs(args, 2, core.MaxArgs); err != nil {
		return core.None, err
	}

	if err := core.ValidateType(args[0], types.Object); err != nil {
		return core.None, err
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

	if err := validateArrayOf(types.String, keys); err != nil {
		return core.None, err
	}

	obj := args[0].(*runtime.Object)
	resultObj := runtime.NewObject()

	var key core.String

	_ = keys.ForEach(ctx, func(c context.Context, keyVal core.Value, idx core.Int) (core.Boolean, error) {
		key = keyVal.(core.String)

		if val, err := obj.Get(c, key); err == nil {
			cloneable, ok := val.(core.Cloneable)

			if ok {
				v, err := cloneable.Clone(c)

				if err != nil {
					return core.False, err
				}

				val = v
			}

			_ = resultObj.Set(c, key, val)
		}

		return true, nil
	})

	return resultObj, nil
}
