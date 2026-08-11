package collections

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// COUNT computes the number of distinct elements in the given collection and returns the count as an integer.
// @param collection {Collection} Collection whose elements are counted.
// @return {Int} Number of elements in the collection.
func Count(ctx context.Context, arg runtime.Value) (runtime.Value, error) {
	collection, err := runtime.CastArg[runtime.Collection](arg, 0)

	if err != nil {
		return runtime.ZeroInt, err
	}

	return collection.Length(ctx)
}
