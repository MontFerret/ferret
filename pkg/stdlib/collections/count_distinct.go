package collections

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/internal/valueset"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// count_distinct computes the number of distinct elements in the given collection and returns the count as an integer.
// @param collection {Collection} Collection whose distinct elements are counted.
// @return {Int} Number of distinct elements in the collection.
func CountDistinct(ctx context.Context, arg runtime.Value) (runtime.Value, error) {
	collection, err := runtime.CastArg[runtime.Collection](arg, 0)

	if err != nil {
		return runtime.ZeroInt, err
	}

	seen := valueset.New(0)

	err = runtime.ForEach(ctx, collection, func(c context.Context, value, idx runtime.Value) (runtime.Boolean, error) {
		if _, err := seen.Add(c, value); err != nil {
			return false, err
		}

		return true, nil
	})

	if err != nil {
		return runtime.ZeroInt, err
	}

	return runtime.Int(seen.Len()), nil
}
