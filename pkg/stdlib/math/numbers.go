package math

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type numberCounts struct {
	total   runtime.Int
	numeric runtime.Int
}

// forEachNumber skips non-numeric elements without coercion for the legacy
// globals. The callback receives numbers and their numeric ordinal. Both counts
// come from traversal so callers never need source length or indexed access.
func forEachNumber(ctx context.Context, source runtime.List, visit func(runtime.Value, runtime.Int)) (numberCounts, error) {
	if err := ctx.Err(); err != nil {
		return numberCounts{}, err
	}

	var counts numberCounts

	err := source.ForEach(ctx, func(c context.Context, value runtime.Value, _ runtime.Int) (runtime.Boolean, error) {
		if err := c.Err(); err != nil {
			return false, err
		}

		counts.total++

		if !runtime.IsNumber(value) {
			return true, nil
		}

		visit(value, counts.numeric)

		counts.numeric++

		return true, nil
	})
	if err != nil {
		return counts, err
	}

	// Preserve the causal traversal error; only successful traversal needs a
	// final cancellation check, including hosts that return without a callback.
	return counts, ctx.Err()
}

func sumNumbers(ctx context.Context, source runtime.List) (float64, numberCounts, error) {
	var sum float64

	counts, err := forEachNumber(ctx, source, func(value runtime.Value, _ runtime.Int) {
		sum += toFloat(value)
	})
	if err != nil {
		return 0, numberCounts{}, err
	}

	return sum, counts, nil
}
