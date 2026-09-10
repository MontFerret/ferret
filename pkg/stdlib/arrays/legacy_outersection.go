package arrays

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/internal/valueset"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// Returns values present in exactly one input array, ignoring duplicates within each input.
// Results follow first-encounter order. Unlike XOR, presence in three inputs is excluded.
// @param arrays {Any[], repeated} At least two arrays.
// @return {Any[]} Values present in exactly one source array.
// @deprecated Use arrays::symmetric_difference for set XOR; it differs for three or more arrays.
func legacyOutersection(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgs(args, 2, runtime.MaxArgs); err != nil {
		return runtime.None, err
	}

	first, err := runtime.CastArgAt[runtime.List](args, 0)
	if err != nil {
		return runtime.None, err
	}

	seen, ordered, err := collectSet(ctx, first)
	if err != nil {
		return runtime.None, err
	}

	repeated := valueset.New(0)
	for index, arg := range args[1:] {
		list, err := runtime.CastArg[runtime.List](arg, index+1)
		if err != nil {
			return runtime.None, err
		}

		source := valueset.New(nativeArrayCapacity(ctx, list))
		err = list.ForEach(ctx, func(ctx context.Context, value runtime.Value, _ runtime.Int) (runtime.Boolean, error) {
			added, err := source.Add(ctx, value)
			if err != nil || !added {
				return runtime.Boolean(err == nil), err
			}

			added, err = seen.Add(ctx, value)
			if err != nil {
				return false, err
			}

			if added {
				_ = ordered.Append(ctx, value)
			} else {
				_, err = repeated.Add(ctx, value)
			}

			return true, err
		})
		if err != nil {
			return runtime.None, err
		}
	}

	return filterSet(ctx, ordered, repeated, false)
}
