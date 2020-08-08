package objects

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

// MERGE_RECURSIVE recursively merge the given objects into a single object.
// @param {Objects, repeated} objects - Objects to merge.
// @return {Object} - Object created by merging.
func MergeRecursive(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, core.MaxArgs)
	if err != nil {
		return values.None, err
	}

	for _, arg := range args {
		if err = core.ValidateType(arg, types.Object); err != nil {
			return values.None, err
		}
	}

	merged := values.NewObject()

	for _, arg := range args {
		merged = merge(merged, arg).(*values.Object)
	}

	return merged.Clone(), nil
}

func merge(src, dst core.Value) core.Value {
	if src.Type() != dst.Type() {
		return dst
	}

	if src.Type() != types.Object {
		return dst
	}

	srcObj := src.(*values.Object)
	dstObj := dst.(*values.Object)

	if dstObj.Length() == 0 {
		return src
	}

	keyObj := values.NewString("")
	exists := values.NewBoolean(false)
	var srcVal core.Value

	dstObj.ForEach(func(val core.Value, key string) bool {
		keyObj = values.NewString(key)

		if srcVal, exists = srcObj.Get(keyObj); exists {
			val = merge(srcVal, val)
		}

		srcObj.Set(keyObj, val)
		return true
	})

	return src
}
