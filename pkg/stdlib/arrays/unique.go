package arrays

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/internal/valueset"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// unique returns all unique elements from a given array.
// @param array {Any[]} Target array.
// @return {Any[]} New array without duplicates.
func Unique(ctx context.Context, arg runtime.Value) (runtime.Value, error) {
	list, err := runtime.CastArg[runtime.List](arg, 0)
	if err != nil {
		return runtime.None, err
	}

	return toUniqueList(ctx, list)
}

func toUniqueList(ctx context.Context, list runtime.List) (runtime.List, error) {
	seen := valueset.New(0)

	return list.Filter(ctx, func(ctx context.Context, value runtime.Value, _ runtime.Int) (runtime.Boolean, error) {
		added, err := seen.Add(ctx, value)

		return runtime.Boolean(added), err
	})
}
