package arrays

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/internal/valueset"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// Intersection returns distinct values present in every input, in first-input order.
// @param arrays {Any[], repeated} At least two arrays.
// @return {Any[]} Common values, retaining their first representatives.
func Intersection(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgs(args, 2, runtime.MaxArgs); err != nil {
		return runtime.None, err
	}

	first, err := runtime.CastArgAt[runtime.List](args, 0)
	if err != nil {
		return runtime.None, err
	}

	candidates, ordered, err := collectSet(ctx, first)
	if err != nil {
		return runtime.None, err
	}

	for index, arg := range args[1:] {
		list, err := runtime.CastArg[runtime.List](arg, index+1)
		if err != nil {
			return runtime.None, err
		}

		next := valueset.New(candidates.Len())
		err = list.ForEach(ctx, func(ctx context.Context, value runtime.Value, _ runtime.Int) (runtime.Boolean, error) {
			found, err := candidates.Contains(ctx, value)
			if err != nil {
				return false, err
			}

			if found {
				_, err = next.Add(ctx, value)
			}

			return true, err
		})
		if err != nil {
			return runtime.None, err
		}

		candidates = next
	}

	return filterSet(ctx, ordered, candidates, true)
}
