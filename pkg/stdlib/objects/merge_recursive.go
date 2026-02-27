package objects

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// MERGE_RECURSIVE recursively merge the given objects into a single object.
// @param {Map, repeated} objects - Maps to merge.
// @return {Map} - Map created by merging.
func MergeRecursive(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgs(args, 1, runtime.MaxArgs); err != nil {
		return runtime.None, err
	}

	if err := runtime.ValidateArgsType(args, runtime.TypeMap); err != nil {
		return runtime.None, err
	}

	first := args[0].(runtime.Map)

	if len(args) == 1 {
		return first.Copy(), nil
	}

	merged, err := first.Empty(ctx)

	if err != nil {
		return runtime.None, err
	}

	for _, arg := range args {
		if err := merge(ctx, merged, arg.(runtime.Map)); err != nil {
			return runtime.None, err
		}
	}

	return merged.Clone(ctx)
}

func merge(ctx context.Context, dst, src runtime.Map) error {
	// If both values are equal, no need to merge
	if runtime.CompareValues(src, dst) == 0 {
		return nil
	}

	return src.ForEach(ctx, func(c context.Context, val, key runtime.Value) (runtime.Boolean, error) {
		switch v := val.(type) {
		case runtime.Map:
			return runtime.True, merge(ctx, v, dst)
		default:
			return runtime.True, dst.Set(c, key, val)
		}
	})
}
