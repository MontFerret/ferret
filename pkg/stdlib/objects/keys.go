package objects

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// keys returns string array of object's keys
// @param obj {Map} The object whose keys you want to extract
// @param sort {Boolean} If sort is true, then the returned keys will be sorted.
// @return {String[]} arrayList that contains object keys.
func Keys(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgs(args, 1, 2); err != nil {
		return runtime.None, err
	}

	if len(args) == 1 {
		return keys1(ctx, args[0])
	}

	return keys2(ctx, args[0], args[1])
}

// keys returns string array of object's keys
// @param obj {Map} The object whose keys you want to extract
// @return {String[]} arrayList that contains object keys.
func keys1(ctx context.Context, arg1 runtime.Value) (runtime.Value, error) {
	return keys2(ctx, arg1, runtime.False)
}

// keys returns string array of object's keys
// @param obj {Map} The object whose keys you want to extract
// @param sort {Boolean} If sort is true, then the returned keys will be sorted.
// @return {String[]} arrayList that contains object keys.
func keys2(ctx context.Context, arg1, arg2 runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgType(arg1, 0, runtime.TypeMap); err != nil {
		return runtime.None, err
	}

	if err := runtime.ValidateArgType(arg2, 1, runtime.TypeBoolean); err != nil {
		return runtime.None, err
	}

	target := arg1.(runtime.Map)
	needSort := arg2.(runtime.Boolean)

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
