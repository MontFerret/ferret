package collections

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime"
)

// COUNT computes the number of distinct elements in the given collection and returns the count as an integer.
func Count(ctx context.Context, arg runtime.Value) (runtime.Value, error) {
	collection, err := runtime.CastCollection(arg)

	if err != nil {
		return runtime.ZeroInt, err
	}

	return collection.Length(ctx)
}
