package math

import (
	"context"
	"errors"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// snapshotNumbers retains the native numeric types in independently owned
// storage. No source copying, sorting, or random access is required.
func snapshotNumbers(ctx context.Context, source runtime.List) ([]runtime.Value, error) {
	var values []runtime.Value

	_, err := forEachNumber(ctx, source, func(value runtime.Value, _ runtime.Int) {
		values = append(values, value)
	})
	if err != nil {
		return nil, err
	}

	return values, nil
}

func sortNumbers(ctx context.Context, values []runtime.Value, ascending bool) error {
	err := runtime.SortSliceWith(ctx, values, func(c context.Context, left, right runtime.Value) (runtime.Ordering, error) {
		if err := c.Err(); err != nil {
			return runtime.Equal, err
		}

		if !ascending {
			left, right = right, left
		}

		return runtime.CompareValues(c, left, right)
	})
	if canceled := ctx.Err(); canceled != nil && !errors.Is(err, canceled) {
		err = errors.Join(err, canceled)
	}

	return err
}
