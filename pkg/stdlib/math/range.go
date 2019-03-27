package math

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

// Range returns an array of numbers in the specified range, optionally with increments other than 1.
// @param start (Int|Float) - The value to start the range at (inclusive).
// @param end (Int|Float) - The value to end the range with (inclusive).
// @param step (Int|Float, optional) - How much to increment in every step, the default is 1.0.
func Range(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 2, 3)

	if err != nil {
		return values.None, err
	}

	err = core.ValidateType(args[0], types.Int, types.Float)

	if err != nil {
		return values.None, err
	}

	err = core.ValidateType(args[1], types.Int, types.Float)

	if err != nil {
		return values.None, err
	}

	var step float64 = 1

	if len(args) > 2 {
		err = core.ValidateType(args[2], types.Int, types.Float)

		if err != nil {
			return values.None, err
		}

		step = toFloat(args[2])
	}

	start := toFloat(args[0])
	end := toFloat(args[1])

	arr := values.NewArray(int(end))

	for i := start; i <= end; i += step {
		arr.Push(values.NewFloat(i))
	}

	return arr, nil
}
