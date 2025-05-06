package objects

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/internal"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

// KEEP_KEYS returns a new object with only given keys.
// @param {hashMap} obj - Source object.
// @param {String, repeated} keys - Keys that need to be kept.
// @return {hashMap} - New hashMap with only given keys.
func KeepKeys(_ context.Context, args ...core.Value) (core.Value, error) {
	if err := core.ValidateArgs(args, 2, core.MaxArgs); err != nil {
		return core.None, err
	}

	if err := core.ValidateType(args[0], types.Object); err != nil {
		return core.None, err
	}

	var keys *internal.Array

	if len(args) == 2 {
		arr, ok := args[1].(*internal.Array)

		if ok {
			keys = arr
		}
	}

	if keys == nil {
		keys = internal.NewArrayWith(args[1:]...)
	}

	if err := validateArrayOf(types.String, keys); err != nil {
		return core.None, err
	}

	obj := args[0].(*internal.Object)
	resultObj := internal.NewObject()

	var key core.String
	var val core.Value
	var exists core.Boolean

	keys.ForEach(func(keyVal core.Value, idx int) bool {
		key = keyVal.(core.String)

		if val, exists = obj.Get(key); exists {
			cloneable, ok := val.(core.Cloneable)

			if ok {
				val = cloneable.Clone()
			}

			resultObj.Set(key, val)
		}

		return true
	})

	return resultObj, nil
}
