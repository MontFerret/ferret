package math

import (
	"context"
	"math"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// clamp limits a number to the inclusive interval from min to max.
// @param value {Int | Float} Input number.
// @param min {Int | Float} Lower bound; must not be NaN or greater than max.
// @param max {Int | Float} Upper bound; must not be NaN. Infinite bounds are allowed.
// @return {Float} The clamped number, using Go math.Min(math.Max(value, min), max) floating-point semantics.
// @throws {TypeError} An argument is not numeric.
// @throws {InvalidArgument} A bound is NaN or min is greater than max.
func Clamp(ctx context.Context, value, min, max runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgValue(value, 0, runtime.AssertNumber); err != nil {
		return runtime.None, err
	}

	if err := runtime.ValidateArgValue(min, 1, runtime.AssertNumber); err != nil {
		return runtime.None, err
	}

	if err := runtime.ValidateArgValue(max, 2, runtime.AssertNumber); err != nil {
		return runtime.None, err
	}

	if number, ok := min.(runtime.Float); ok && math.IsNaN(float64(number)) {
		return runtime.None, runtime.ArgError(runtime.Error(runtime.ErrInvalidArgument, "minimum must not be NaN"), 1)
	}

	if number, ok := max.(runtime.Float); ok && math.IsNaN(float64(number)) {
		return runtime.None, runtime.ArgError(runtime.Error(runtime.ErrInvalidArgument, "maximum must not be NaN"), 2)
	}

	// Check the original bounds before float conversion can lose integer precision.
	order, err := runtime.CompareValues(ctx, min, max)
	if err != nil {
		return runtime.None, err
	}

	if order == runtime.Greater {
		return runtime.None, runtime.ArgError(runtime.Error(runtime.ErrInvalidArgument, "minimum must not exceed maximum"), 1)
	}

	return runtime.NewFloat(math.Min(math.Max(toFloat(value), toFloat(min)), toFloat(max))), nil
}
