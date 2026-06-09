package collections

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/internal/valueset"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// COUNT_DISTINCT computes the number of distinct elements in the given collection and returns the count as an integer.
func CountDistinct(ctx context.Context, arg runtime.Value) (runtime.Value, error) {
	collection, err := runtime.CastArg[runtime.Collection](arg, 0)

	if err != nil {
		return runtime.ZeroInt, err
	}

	seen := valueset.New(0)

	err = runtime.ForEach(ctx, collection, func(c context.Context, value, idx runtime.Value) (runtime.Boolean, error) {
		seen.Add(value)

		return true, nil
	})

	if err != nil {
		return runtime.ZeroInt, err
	}

	return runtime.Int(seen.Len()), nil
}
