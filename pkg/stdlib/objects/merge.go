package objects

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime"
)

// MERGE merge the given objects into a single object.
// @param {hashMap, repeated} objects - Objects to merge.
// @return {hashMap} - hashMap created by merging.
// TODO: REWRITE TO USE LIST & MAP instead
func Merge(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgs(args, 1, runtime.MaxArgs); err != nil {
		return runtime.None, err
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

	if err := validateArrayOf(ctx, runtime.TypeObject, objs); err != nil {
		return runtime.None, err
	}

	return mergeArray(ctx, objs)
}

func mergeArray(ctx context.Context, arr *runtime.Array) (*runtime.Object, error) {
	merged, obj := runtime.NewObject(), runtime.NewObject()

	_ = arr.ForEach(ctx, func(c context.Context, arrValue runtime.Value, arrIdx runtime.Int) (runtime.Boolean, error) {
		obj = arrValue.(*runtime.Object)

		_ = obj.ForEach(c, func(_ context.Context, objValue, objKey runtime.Value) (runtime.Boolean, error) {
			cloneable, ok := objValue.(runtime.Cloneable)

			if ok {
				clone, err := cloneable.Clone(c)

				if err != nil {
					return runtime.False, err
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

func validateArrayOf(ctx context.Context, typ runtime.Type, arr *runtime.Array) error {
	size, err := arr.Length(ctx)

	if err != nil {
		return err
	}

	for idx := runtime.ZeroInt; idx < size; idx++ {
		item, err := arr.Get(ctx, idx)

		if err != nil {
			return err
		}

		if err := runtime.ValidateType(item, typ); err != nil {
			return err
		}
	}

	return nil
}
