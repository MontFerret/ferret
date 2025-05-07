package math

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime"

	"github.com/MontFerret/ferret/pkg/runtime/core"
)

// RANGE returns an array of numbers in the specified range, optionally with increments other than 1.
// @param {Int | Float} start - The value to start the range at (inclusive).
// @param {Int | Float} end - The value to end the range with (inclusive).
// @param {Int | Float} [step=1.0] - How much to increment in every step.
// @return {Int[] | Float[]} - arrayList of numbers in the specified range, optionally with increments other than 1.
func Range(ctx context.Context, args ...core.Value) (core.Value, error) {
	if err := core.ValidateArgs(args, 2, 3); err != nil {
		return core.None, err
	}

	if err := runtime.ValidateArgType(args, 0, runtime.AssertNumber); err != nil {
		return core.None, err
	}

	if err := runtime.ValidateArgType(args, 1, runtime.AssertNumber); err != nil {
		return core.None, err
	}

	var step float64 = 1

	if len(args) > 2 {
		if err := runtime.ValidateArgType(args, 2, runtime.AssertNumber); err != nil {
			return core.None, err
		}

		step = toFloat(args[2])
	}

	start := toFloat(args[0])
	end := toFloat(args[1])

	arr := runtime.NewArray(int(end))

	for i := start; i <= end; i += step {
		_ = arr.Add(ctx, core.NewFloat(i))
	}

	return arr, nil
}
