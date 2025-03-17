package objects

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/internal"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

// MERGE merge the given objects into a single object.
// @param {hashMap, repeated} objects - Objects to merge.
// @return {hashMap} - hashMap created by merging.
func Merge(_ context.Context, args ...core.Value) (core.Value, error) {
	if err := core.ValidateArgs(args, 1, core.MaxArgs); err != nil {
		return core.None, err
	}

	var objs *internal.Array

	if len(args) == 1 {
		arr, ok := args[0].(*internal.Array)

		if ok {
			objs = arr
		}
	}

	if objs == nil {
		objs = internal.NewArrayWith(args...)
	}

	if err := validateArrayOf(types.Object, objs); err != nil {
		return core.None, err
	}

	return mergeArray(objs), nil
}

func mergeArray(arr *internal.Array) *internal.Object {
	merged, obj := internal.NewObject(), internal.NewObject()

	arr.ForEach(func(arrValue core.Value, arrIdx int) bool {
		obj = arrValue.(*internal.Object)
		obj.ForEach(func(objValue core.Value, objKey string) bool {
			cloneable, ok := objValue.(core.Cloneable)

			if ok {
				objValue = cloneable.Clone()
			}

			merged.Set(core.NewString(objKey), objValue)

			return true
		})
		return true
	})

	return merged
}

func validateArrayOf(typ core.Type, arr *internal.Array) (err error) {
	for idx := 0; idx < arr.Length(); idx++ {
		if err != nil {
			break
		}
		err = core.ValidateType(arr.Get(idx), typ)
	}
	return
}
