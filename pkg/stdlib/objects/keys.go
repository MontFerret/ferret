package objects

import (
	"context"
	"sort"

	"github.com/MontFerret/ferret/pkg/runtime"
)

// KEYS returns string array of object's keys
// @param {hashMap} obj - The object whose keys you want to extract
// @param {Boolean} [sort=False] - If sort is true, then the returned keys will be sorted.
// @return {String[]} - arrayList that contains object keys.
// TODO: REWRITE TO USE LIST & MAP instead
func Keys(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
	err := runtime.ValidateArgs(args, 1, 2)
	if err != nil {
		return runtime.None, err
	}

	err = runtime.ValidateType(args[0], runtime.TypeObject)
	if err != nil {
		return runtime.None, err
	}

	obj := args[0].(*runtime.Object)
	needSort := false

	if len(args) == 2 {
		err = runtime.ValidateType(args[1], runtime.TypeBoolean)
		if err != nil {
			return runtime.None, err
		}

		needSort = bool(args[1].(runtime.Boolean))
	}

	size, err := obj.Length(ctx)

	if err != nil {
		return runtime.None, err
	}

	oKeys := make([]string, 0, size)

	_ = obj.ForEach(ctx, func(c context.Context, value, key runtime.Value) (runtime.Boolean, error) {
		oKeys = append(oKeys, key.String())

		return true, nil
	})

	keys := sort.StringSlice(oKeys)
	keysArray := runtime.NewArray(len(keys))

	if needSort {
		keys.Sort()
	}

	for _, key := range keys {
		_ = keysArray.Add(ctx, runtime.NewString(key))
	}

	return keysArray, nil
}
