package math

import (
	"context"
	"fmt"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type (
	numberPolicy uint8

	numberCounts struct {
		total   runtime.Int
		numeric runtime.Int
	}
)

const (
	filterNonNumbers numberPolicy = iota
	strictNumbers
)

// forEachNumber either rejects or skips non-numeric elements without coercion.
// The callback receives numbers and their numeric ordinal. Both counts
// come from traversal so callers never need source length or indexed access.
func forEachNumber(ctx context.Context, source runtime.List, policy numberPolicy, visit func(runtime.Value, runtime.Int)) (numberCounts, error) {
	var counts numberCounts

	err := source.ForEach(ctx, func(_ context.Context, value runtime.Value, index runtime.Int) (runtime.Boolean, error) {
		counts.total++

		if !runtime.IsNumber(value) {
			if policy == strictNumbers {
				return false, numberElementError(value, index)
			}

			return true, nil
		}

		visit(value, counts.numeric)

		counts.numeric++

		return true, nil
	})

	return counts, err
}

func sumNumbers(ctx context.Context, source runtime.List, policy numberPolicy) (float64, numberCounts, error) {
	var sum float64

	counts, err := forEachNumber(ctx, source, policy, func(value runtime.Value, _ runtime.Int) {
		sum += toFloat(value)
	})
	if err != nil {
		return 0, numberCounts{}, err
	}

	return sum, counts, nil
}

func numberElementError(value runtime.Value, index runtime.Int) error {
	return runtime.ArgError(fmt.Errorf("element at index %d: %w", index, runtime.AssertNumber(value)), 0)
}
