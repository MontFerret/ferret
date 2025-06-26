package objects

import (
	"context"
	"sort"

	"github.com/MontFerret/ferret/pkg/runtime"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

// KEYS returns string array of object's keys
// @param {hashMap} obj - The object whose keys you want to extract
// @param {Boolean} [sort=False] - If sort is true, then the returned keys will be sorted.
// @return {String[]} - arrayList that contains object keys.
// TODO: REWRITE TO USE LIST & MAP instead
func Keys(ctx context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, 2)
	if err != nil {
		return core.None, err
	}

	err = core.ValidateType(args[0], types.Object)
	if err != nil {
		return core.None, err
	}

	obj := args[0].(*runtime.Object)
	needSort := false

	if len(args) == 2 {
		err = core.ValidateType(args[1], types.Boolean)
		if err != nil {
			return core.None, err
		}

		needSort = bool(args[1].(core.Boolean))
	}

	size, err := obj.Length(ctx)

	if err != nil {
		return core.None, err
	}

	oKeys := make([]string, 0, size)

	_ = obj.ForEach(ctx, func(c context.Context, value, key core.Value) (core.Boolean, error) {
		oKeys = append(oKeys, key.String())

		return true, nil
	})

	keys := sort.StringSlice(oKeys)
	keysArray := runtime.NewArray(len(keys))

	if needSort {
		keys.Sort()
	}

	for _, key := range keys {
		_ = keysArray.Add(ctx, core.NewString(key))
	}

	return keysArray, nil
}
