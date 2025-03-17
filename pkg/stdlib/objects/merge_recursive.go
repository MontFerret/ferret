package objects

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/internal"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

// MERGE_RECURSIVE recursively merge the given objects into a single object.
// @param {Objects, repeated} objects - Objects to merge.
// @return {hashMap} - hashMap created by merging.
func MergeRecursive(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, core.MaxArgs)
	if err != nil {
		return core.None, err
	}

	for _, arg := range args {
		if err = core.ValidateType(arg, types.Object); err != nil {
			return core.None, err
		}
	}

	merged := internal.NewObject()

	for _, arg := range args {
		merged = merge(merged, arg).(*internal.Object)
	}

	return merged.Clone(), nil
}

func merge(src, dst core.Value) core.Value {
	if !types.Equal(core.Reflect(src), core.Reflect(dst)) {
		return src
	}

	srcObj, ok := src.(*internal.Object)

	if !ok {
		return dst
	}

	dstObj, ok := dst.(*internal.Object)

	if !ok {
		return src
	}

	if dstObj.Length() == 0 {
		return src
	}

	keyObj := core.NewString("")
	exists := core.NewBoolean(false)
	var srcVal core.Value

	dstObj.ForEach(func(val core.Value, key string) bool {
		keyObj = core.NewString(key)

		if srcVal, exists = srcObj.Get(keyObj); exists {
			val = merge(srcVal, val)
		}

		srcObj.Set(keyObj, val)
		return true
	})

	return src
}
