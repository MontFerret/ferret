package math

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

// RAND return a pseudo-random number between 0 and 1.
// @param {Int | Float} [max] - Upper limit.
// @param {Int | Float} [min] - Lower limit.
// @return {Float} - A number greater than 0 and less than 1.
func Rand(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 0, 2)

	if err != nil {
		return values.None, err
	}

	if len(args) == 0 {
		return values.NewFloat(core.RandomDefault()), nil
	}

	var maxValue float64
	var minValue float64

	maxValue = float64(values.ToFloat(args[0]))

	if len(args) > 1 {
		minValue = float64(values.ToFloat(args[1]))
	} else {
		maxValue, minValue = core.NumberBoundaries(maxValue)
	}

	return values.NewFloat(core.Random(maxValue, minValue)), nil
}
