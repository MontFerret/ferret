package math

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

// AVERAGE Returns the average (arithmetic mean) of the values in array.
// @param {Int[] | Float[]} array - Array of numbers.
// @return {Float} - The average of the values in array.
func Average(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, 1)

	if err != nil {
		return values.None, err
	}

	err = core.ValidateType(args[0], types.Array)

	if err != nil {
		return values.None, err
	}

	arr := args[0].(*values.Array)

	if arr.Length() == 0 {
		return values.None, nil
	}

	var sum float64

	arr.ForEach(func(value core.Value, idx int) bool {
		err = core.ValidateType(value, types.Float, types.Int)

		if err != nil {
			return false
		}

		sum += toFloat(value)

		return true
	})

	if err != nil {
		return values.None, nil
	}

	count := arr.Length()

	return values.Float(sum / float64(count)), nil
}
