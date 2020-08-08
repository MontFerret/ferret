package objects

import (
	"context"
	"sort"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

// KEYS returns string array of object's keys
// @param {Object} obj - The object whose keys you want to extract
// @param {Boolean} [sort=False] - If sort is true, then the returned keys will be sorted.
// @return {String[]} - Array that contains object keys.
func Keys(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, 2)
	if err != nil {
		return values.None, err
	}

	err = core.ValidateType(args[0], types.Object)
	if err != nil {
		return values.None, err
	}

	obj := args[0].(*values.Object)
	needSort := false

	if len(args) == 2 {
		err = core.ValidateType(args[1], types.Boolean)
		if err != nil {
			return values.None, err
		}

		needSort = bool(args[1].(values.Boolean))
	}

	oKeys := make([]string, 0, obj.Length())

	obj.ForEach(func(value core.Value, key string) bool {
		oKeys = append(oKeys, key)

		return true
	})

	keys := sort.StringSlice(oKeys)
	keysArray := values.NewArray(len(keys))

	if needSort {
		keys.Sort()
	}

	for _, key := range keys {
		keysArray.Push(values.NewString(key))
	}

	return keysArray, nil
}
