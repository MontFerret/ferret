package objects

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime"
)

// MERGE_RECURSIVE recursively merge the given objects into a single object.
// @param {Objects, repeated} objects - Objects to merge.
// @return {hashMap} - hashMap created by merging.
// TODO: REWRITE TO USE LIST & MAP instead
func MergeRecursive(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
	err := runtime.ValidateArgs(args, 1, runtime.MaxArgs)
	if err != nil {
		return runtime.None, err
	}

	for _, arg := range args {
		if err = runtime.ValidateType(arg, runtime.TypeObject); err != nil {
			return runtime.None, err
		}
	}

	merged := runtime.NewObject()

	for _, arg := range args {
		out, err := merge(ctx, merged, arg)

		if err != nil {
			return runtime.None, err
		}

		merged = out.(*runtime.Object)
	}

	return merged.Clone(ctx)
}

func merge(ctx context.Context, src, dst runtime.Value) (runtime.Value, error) {
	if runtime.CompareValues(src, dst) != 0 {
		return src, nil
	}

	srcObj, ok := src.(*runtime.Object)

	if !ok {
		return dst, nil
	}

	dstObj, ok := dst.(*runtime.Object)

	if !ok {
		return src, nil
	}

	size, err := dstObj.Length(ctx)

	if err != nil {
		return runtime.None, err
	}

	if size == 0 {
		return src, nil
	}

	var srcVal runtime.Value

	_ = dstObj.ForEach(ctx, func(c context.Context, val, key runtime.Value) (runtime.Boolean, error) {
		if srcVal, err = srcObj.Get(c, key); err == nil {
			v, err := merge(ctx, srcVal, val)

			if err != nil {
				return runtime.False, err
			}

			val = v
		}

		_ = srcObj.Set(c, key, val)
		return true, nil
	})

	return src, nil
}
