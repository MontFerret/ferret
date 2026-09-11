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

// MergeMutable merges sources into the original target and returns that same map.
// Incoming values are cloned or copied; later sources win. Errors may leave
// earlier updates applied. A target without sources is returned unchanged.
// @param target {Map} Map to mutate.
// @param sources {Map|Map[], repeated} Zero or more maps, or one list of maps.
// @return {Map} The original mutated target.
func MergeMutable(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
	return mergeMutableObjects(ctx, args, runtime.MergeMapsInto)
}

func mergeObjects(ctx context.Context, args []runtime.Value, mergeInto func(context.Context, runtime.Map, ...runtime.Map) error) (runtime.Value, error) {
	if err := runtime.ValidateArgs(args, 1, runtime.MaxArgs); err != nil {
		return runtime.None, err
	}

	sources, listForm, err := normalizeMergeArgs(ctx, args, 0)
	if err != nil {
		return runtime.None, err
	}

	if len(sources) == 0 {
		return runtime.NewObject(), nil
	}

	dst, err := sources[0].New(ctx)
	if err != nil {
		return runtime.None, mergeArgumentError(err, 0, listForm, 0)
	}

	for index, src := range sources {
		if err := mergeInto(ctx, dst, src); err != nil {
			return runtime.None, mergeArgumentError(err, index, listForm, 0)
		}
	}

	return dst, nil
}

func mergeMutableObjects(ctx context.Context, args []runtime.Value, mergeInto func(context.Context, runtime.Map, ...runtime.Map) error) (runtime.Value, error) {
	if err := runtime.ValidateArgs(args, 1, runtime.MaxArgs); err != nil {
		return runtime.None, err
	}

	if err := ctx.Err(); err != nil {
		return runtime.None, err
	}

	dst, err := runtime.CastArg[runtime.Map](args[0], 0)
	if err != nil {
		return runtime.None, err
	}

	sources, listForm, err := normalizeMergeArgs(ctx, args[1:], 1)
	if err != nil {
		return runtime.None, err
	}

	for index, src := range sources {
		if err := mergeInto(ctx, dst, src); err != nil {
			return runtime.None, mergeArgumentError(err, index, listForm, 1)
		}
	}

	return dst, nil
}
