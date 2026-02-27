package objects

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// KEYS returns string array of object's keys
// @param {Map} obj - The object whose keys you want to extract
// @param {Boolean} [sort=False] - If sort is true, then the returned keys will be sorted.
// @return {String[]} - arrayList that contains object keys.
func Keys(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgs(args, 1, 2); err != nil {
		return runtime.None, err
	}

	if err := runtime.ValidateArgTypeAt(args, 0, runtime.TypeMap); err != nil {
		return runtime.None, err
	}

	target := args[0].(runtime.Map)
	needSort := runtime.False

	if len(args) == 2 {
		if err := runtime.ValidateArgTypeAt(args, 1, runtime.TypeBoolean); err != nil {
			return runtime.None, err
		}

		needSort = args[1].(runtime.Boolean)
	}

	keys, err := target.Keys(ctx)

	if err != nil {
		return runtime.None, err
	}

	if needSort {
		if err := keys.SortAsc(ctx); err != nil {
			return runtime.None, err
		}
	}

	return keys, nil
}
