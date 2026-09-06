package objects

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// Merge copies maps into an independent destination; later values replace earlier values.
// @param values {Map|Map[], repeated} Variadic maps, or one list of maps.
// @return {Map} Shallow merge with cloned or copied values.
func Merge(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
	return mergeObjects(ctx, args, runtime.MergeMapsInto)
}

func mergeObjects(ctx context.Context, args []runtime.Value, mergeInto func(context.Context, runtime.Map, ...runtime.Map) error) (runtime.Value, error) {
	sources, listForm, err := normalizeMergeArgs(ctx, args)
	if err != nil {
		return runtime.None, err
	}

	if len(sources) == 0 {
		return runtime.NewObject(), nil
	}

	dst, err := sources[0].Empty(ctx)
	if err != nil {
		return runtime.None, mergeArgumentError(err, 0, listForm)
	}

	for index, src := range sources {
		if err := mergeInto(ctx, dst, src); err != nil {
			return runtime.None, mergeArgumentError(err, index, listForm)
		}
	}

	return dst, nil
}
