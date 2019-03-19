package math

import (
	"context"
	"math/rand"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

// Rand return a pseudo-random number between 0 and 1.
// @param max (Float|Int, optional) - Upper limit.
// @param min (Float|Int, optional) - Lower limit.
// @returns (Float) - A number greater than 0 and less than 1.
func Rand(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 0, 2)

	if err != nil {
		return values.None, err
	}

	if len(args) == 0 {
		return values.NewFloat(rand.Float64()), nil
	}

	var max float64
	var min float64

	arg1, err := values.ToFloat(args[0])

	if err != nil {
		return values.None, err
	}

	max = float64(arg1)

	if len(args) > 1 {
		arg2, err := values.ToFloat(args[1])

		if err != nil {
			return values.None, err
		}

		min = float64(arg2)
	} else {
		max, min = core.RandomBoundaries(max)
	}

	return values.NewFloat(core.Random(max, min)), nil
}
