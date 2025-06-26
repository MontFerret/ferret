package objects

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

// MERGE merge the given objects into a single object.
// @param {hashMap, repeated} objects - Objects to merge.
// @return {hashMap} - hashMap created by merging.
// TODO: REWRITE TO USE LIST & MAP instead
func Merge(ctx context.Context, args ...core.Value) (core.Value, error) {
	if err := core.ValidateArgs(args, 1, core.MaxArgs); err != nil {
		return core.None, err
	}

	var objs *runtime.Array

	if len(args) == 1 {
		arr, ok := args[0].(*runtime.Array)

		if ok {
			objs = arr
		}
	}

	if objs == nil {
		objs = runtime.NewArrayWith(args...)
	}

	if err := validateArrayOf(ctx, types.Object, objs); err != nil {
		return core.None, err
	}

	return mergeArray(ctx, objs)
}

func mergeArray(ctx context.Context, arr *runtime.Array) (*runtime.Object, error) {
	merged, obj := runtime.NewObject(), runtime.NewObject()

	_ = arr.ForEach(ctx, func(c context.Context, arrValue core.Value, arrIdx core.Int) (core.Boolean, error) {
		obj = arrValue.(*runtime.Object)

		_ = obj.ForEach(c, func(_ context.Context, objValue, objKey core.Value) (core.Boolean, error) {
			cloneable, ok := objValue.(core.Cloneable)

			if ok {
				clone, err := cloneable.Clone(c)

				if err != nil {
					return core.False, err
				}

				objValue = clone
			}

			_ = merged.Set(c, objKey, objValue)

			return true, nil
		})

		return true, nil
	})

	return merged, nil
}

func validateArrayOf(ctx context.Context, typ core.Type, arr *runtime.Array) error {
	size, err := arr.Length(ctx)

	if err != nil {
		return err
	}

	for idx := runtime.ZeroInt; idx < size; idx++ {
		item, err := arr.Get(ctx, idx)

		if err != nil {
			return err
		}

		if err := core.ValidateType(item, typ); err != nil {
			return err
		}
	}

	return nil
}
