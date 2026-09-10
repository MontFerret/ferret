package arrays

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/internal/valueset"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// SymmetricDifference returns values present in an odd number of input arrays.
// Duplicates within one input count once. Results follow first-encounter order.
// @param arrays {Any[], repeated} At least two arrays.
// @return {Any[]} Odd-presence values, retaining their first representatives.
func SymmetricDifference(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgs(args, 2, runtime.MaxArgs); err != nil {
		return runtime.None, err
	}

	seen := valueset.New(0)
	present := valueset.New(0)
	ordered := runtime.NewArray(0)
	for index, arg := range args {
		list, err := runtime.CastArg[runtime.List](arg, index)
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
			}

			removed, err := present.Remove(ctx, value)
			if err != nil {
				return false, err
			}

			if !removed {
				_, err = present.Add(ctx, value)
			}

			return true, err
		})
		if err != nil {
			return runtime.None, err
		}
	}

	return filterSet(ctx, ordered, present, true)
}
