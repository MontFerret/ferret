package math

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
	"math"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

// Median returns the median of the values in array.
// @param array (Array) - Array of numbers.
// @returns (Float) - The median of the values in array.
func Median(_ context.Context, args ...core.Value) (core.Value, error) {
	var err error
	err = core.ValidateArgs(args, 1, 1)

	if err != nil {
		return values.None, err
	}

	err = core.ValidateType(args[0], types.Array)

	if err != nil {
		return values.None, err
	}

	arr := args[0].(*values.Array)
	sorted := arr.Sort()

	l := sorted.Length()

	var median core.Value

	if l == 0 {
		return values.NewFloat(math.NaN()), nil
	} else if l%2 == 0 {
		median, err = mean(sorted.Slice(l/2-1, l/2+1))

		if err != nil {
			return values.None, nil
		}
	} else {
		median = sorted.Get(l / 2)
	}

	if err != nil {
		return values.None, nil
	}

	return median, nil
}
