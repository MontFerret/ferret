package math

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime"
)

// RAND return a pseudo-random number between 0 and 1.
// @param {Int | Float} [max] - Upper limit.
// @param {Int | Float} [min] - Lower limit.
// @return {Float} - A number greater than 0 and less than 1.
func Rand(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgs(args, 0, 2); err != nil {
		return runtime.None, err
	}

	if len(args) == 0 {
		return runtime.NewFloat(runtime.RandomDefault()), nil
	}

	arg1, err := runtime.ToFloat(ctx, args[0])

	if err != nil {
		return runtime.None, err
	}

	var max float64
	var min float64

	max = float64(arg1)

	if len(args) > 1 {
		arg2, err := runtime.ToFloat(ctx, args[1])

		if err != nil {
			return runtime.None, err
		}

		min = float64(arg2)
	} else {
		max, min = runtime.NumberBoundaries(max)
	}

	return runtime.NewFloat(runtime.Random(max, min)), nil
}
