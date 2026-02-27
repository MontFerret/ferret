package objects

import (
	"context"
	"errors"

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
		cloned, err := runtime.CloneOrCopy(ctx, first)
		if err != nil {
			return runtime.None, err
		}

		return cloned, nil
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
			existing, err := dst.Get(c, key)
			if err == nil {
				if existingMap, ok := existing.(runtime.Map); ok {
					if err := merge(c, existingMap, v); err != nil {
						return false, err
					}

					return true, nil
				}
			} else if !errors.Is(err, runtime.ErrNotFound) {
				return false, err
			}

			cloned, err := v.Clone(c)
			if err != nil {
				return false, err
			}

			return true, dst.Set(c, key, cloned.(runtime.Map))
		default:
			return true, dst.Set(c, key, val)
		}
	})
}
