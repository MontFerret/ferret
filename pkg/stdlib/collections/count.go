package collections

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/runtime/core"
)

// COUNT computes the number of distinct elements in the given collection and returns the count as an integer.
func Count(ctx context.Context, arg core.Value) (core.Value, error) {
	collection, err := runtime.CastCollection(arg)

	if err != nil {
		return runtime.ZeroInt, err
	}

	return collection.Length(ctx)
}
