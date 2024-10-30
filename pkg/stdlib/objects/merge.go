package objects

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

// MERGE merge the given objects into a single object.
// @param {Object, repeated} objects - Objects to merge.
// @return {Object} - Object created by merging.
func Merge(_ context.Context, args ...core.Value) (core.Value, error) {
	if err := core.ValidateArgs(args, 1, core.MaxArgs); err != nil {
		return values.None, err
	}

	var objs *values.Array

	if len(args) == 1 {
		arr, ok := args[0].(*values.Array)

		if ok {
			objs = arr
		}
	}

	if objs == nil {
		objs = values.NewArrayWith(args...)
	}

	if err := validateArrayOf(types.Object, objs); err != nil {
		return values.None, err
	}

	return mergeArray(objs), nil
}

func mergeArray(arr *values.Array) *values.Object {
	merged, obj := values.NewObject(), values.NewObject()

	arr.ForEach(func(arrValue core.Value, arrIdx int) bool {
		obj = arrValue.(*values.Object)
		obj.ForEach(func(objValue core.Value, objKey string) bool {
			cloneable, ok := objValue.(core.Cloneable)

			if ok {
				objValue = cloneable.Clone()
			}

			merged.Set(values.NewString(objKey), objValue)

			return true
		})
		return true
	})

	return merged
}

func validateArrayOf(typ core.Type, arr *values.Array) (err error) {
	for idx := 0; idx < arr.Length(); idx++ {
		if err != nil {
			break
		}
		err = core.ValidateType(arr.Get(idx), typ)
	}
	return
}
