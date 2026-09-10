package arrays

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/internal/valueset"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// RemoveAny removes all occurrences of values in the second array.
// Other duplicates and the original element order are preserved.
// @param array {Any[]} Source array.
// @param values {Any[]} Values whose occurrences are removed.
// @return {Any[]} Remaining elements in their original order.
func RemoveAny(ctx context.Context, array, values runtime.Value) (runtime.Value, error) {
	list, targets, err := runtime.CastArgs2[runtime.List, runtime.List](array, values)
	if err != nil {
		return runtime.None, err
	}

	excluded := valueset.New(0)
	err = targets.ForEach(ctx, func(ctx context.Context, value runtime.Value, _ runtime.Int) (runtime.Boolean, error) {
		_, err := excluded.Add(ctx, value)

		return true, err
	})
	if err != nil {
		return runtime.None, err
	}

	return list.Filter(ctx, func(ctx context.Context, value runtime.Value, _ runtime.Int) (runtime.Boolean, error) {
		found, err := excluded.Contains(ctx, value)

		return runtime.Boolean(!found), err
	})
}
