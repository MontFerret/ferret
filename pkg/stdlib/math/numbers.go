package math

import (
	"context"
	"errors"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// forEachNumber counts successfully visited numbers without requiring list
// length or indexed access. Validation never coerces or skips an element.
func forEachNumber(ctx context.Context, source runtime.List, visit func(runtime.Value, runtime.Int)) (runtime.Int, error) {
	if err := ctx.Err(); err != nil {
		return 0, err
	}

	var count runtime.Int

	err := source.ForEach(ctx, func(c context.Context, value runtime.Value, index runtime.Int) (runtime.Boolean, error) {
		if err := c.Err(); err != nil {
			return false, err
		}

		if err := runtime.AssertNumber(value); err != nil {
			return false, runtime.ArgError(runtime.Errorf(err, "at index %d", index), 0)
		}

		visit(value, count)

		count++

		return true, nil
	})
	if canceled := ctx.Err(); canceled != nil && !errors.Is(err, canceled) {
		err = errors.Join(err, canceled)
	}

	return count, err
}

func sumNumbers(ctx context.Context, source runtime.List) (float64, runtime.Int, error) {
	var sum float64

	count, err := forEachNumber(ctx, source, func(value runtime.Value, _ runtime.Int) {
		sum += toFloat(value)
	})
	if err != nil {
		return 0, 0, err
	}

	return sum, count, nil
}
