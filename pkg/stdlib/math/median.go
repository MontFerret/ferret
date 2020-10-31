package math

import (
	"context"
	"math"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

// MEDIAN returns the median of the values in array.
// @param {Int[] | Float[]} array - Array of numbers.
// @return {Float} - The median of the values in array.
func Median(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, 1)

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

	switch {
	case l == 0:
		return values.NewFloat(math.NaN()), nil
	case l%2 == 0:
		median, err = mean(sorted.Slice(l/2-1, l/2+1))

		if err != nil {
			return values.None, nil
		}
	default:
		median = sorted.Get(l / 2)
	}

	return median, nil
}
