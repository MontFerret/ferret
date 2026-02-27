package objects

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// MERGE merge the given objects into a single object.
// @param {Map, repeated} objects - Maps to merge.
// @return {Map} - Map created by merging.
func Merge(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgs(args, 1, runtime.MaxArgs); err != nil {
		return runtime.None, err
	}

	var objs runtime.List

	if len(args) == 1 {
		arr, ok := args[0].(runtime.List)

		if ok {
			objs = arr
		}
	}

	if objs == nil {
		objs = runtime.NewArrayWith(args...)
	}

	if err := runtime.AssertItemsOf(ctx, objs, runtime.AssertMap); err != nil {
		return runtime.None, err
	}

	merged, err := mergeMaps(ctx, objs)
	if err != nil {
		return runtime.None, err
	}

	cloned, err := runtime.CloneOrCopy(ctx, merged)
	if err != nil {
		return runtime.None, err
	}

	return cloned, nil
}

func mergeMaps(ctx context.Context, arr runtime.List) (runtime.Map, error) {
	var obj runtime.Map

	first, err := arr.First(ctx)

	if err != nil {
		return nil, err
	}

	if first == nil || first == runtime.None {
		return runtime.NewObject(), nil
	}

	firstMap := first.(runtime.Map)
	merged, err := firstMap.Empty(ctx)

	if err != nil {
		return runtime.NewObject(), err
	}

	return merged, arr.ForEach(ctx, func(c context.Context, arrValue runtime.Value, arrIdx runtime.Int) (runtime.Boolean, error) {
		obj = arrValue.(runtime.Map)

		return true, merged.Merge(ctx, obj)
	})
}
