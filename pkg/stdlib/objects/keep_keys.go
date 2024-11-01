package objects

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

// KEEP_KEYS returns a new object with only given keys.
// @param {Object} obj - Source object.
// @param {String, repeated} keys - Keys that need to be kept.
// @return {Object} - New Object with only given keys.
func KeepKeys(_ context.Context, args ...core.Value) (core.Value, error) {
	if err := core.ValidateArgs(args, 2, core.MaxArgs); err != nil {
		return values.None, err
	}

	if err := core.ValidateType(args[0], types.Object); err != nil {
		return values.None, err
	}

	var keys *values.Array

	if len(args) == 2 {
		arr, ok := args[1].(*values.Array)

		if ok {
			keys = arr
		}
	}

	if keys == nil {
		keys = values.NewArrayWith(args[1:]...)
	}

	if err := validateArrayOf(types.String, keys); err != nil {
		return values.None, err
	}

	obj := args[0].(*values.Object)
	resultObj := values.NewObject()

	var key values.String
	var val core.Value
	var exists values.Boolean

	keys.ForEach(func(keyVal core.Value, idx int) bool {
		key = keyVal.(values.String)

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
