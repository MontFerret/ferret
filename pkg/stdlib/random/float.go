package random

import (
	"context"
	"math"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// Float draws a pseudo-random Float with zero or two numeric bounds. Ranged
// draws use [min, max); equal bounds return Float(min) without advancing the RNG.
// Non-finite bounds and unequal intervals containing no Float are invalid.
func Float(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
	if len(args) == 0 {
		return float0(ctx)
	}

	if err := runtime.ValidateArgs(args, 2, 2); err != nil {
		return runtime.None, err
	}

	return float2(ctx, args[0], args[1])
}

// float draws a non-cryptographic pseudo-random number from the Session source.
// @return {Float} A number in the half-open interval [0, 1).
func float0(ctx context.Context) (runtime.Value, error) {
	src, err := source(ctx)
	if err != nil {
		return runtime.None, err
	}

	return runtime.Float(src.Float64()), nil
}

// float draws a non-cryptographic pseudo-random number between finite bounds.
// @param min {Int | Float} Inclusive lower bound; must be finite and not exceed max.
// @param max {Int | Float} Exclusive upper bound; must be finite.
// @return {Float} A number in [min, max), or Float(min) for equal bounds without drawing. Original numeric bounds are preserved for unequal intervals.
// @throws {TypeError} A bound is not a native Int or Float.
// @throws {InvalidArgument} Bounds are non-finite, reversed, or unequal with no representable Float in the interval.
func float2(ctx context.Context, min, max runtime.Value) (runtime.Value, error) {
	for pos, bound := range []runtime.Value{min, max} {
		if err := runtime.ValidateArgValue(bound, pos, runtime.AssertNumber); err != nil {
			return runtime.None, err
		}

		if floating, ok := bound.(runtime.Float); ok && (math.IsNaN(float64(floating)) || math.IsInf(float64(floating), 0)) {
			return runtime.None, runtime.ArgError(runtime.Error(runtime.ErrInvalidArgument, "bound must be finite"), pos)
		}
	}

	order, err := runtime.CompareValues(ctx, min, max)
	if err != nil {
		return runtime.None, err
	}

	if order == runtime.Greater {
		return runtime.None, runtime.ArgError(runtime.Error(runtime.ErrInvalidArgument, "minimum must not exceed maximum"), 0)
	}

	if order == runtime.Equal {
		if _, err := source(ctx); err != nil {
			return runtime.None, err
		}

		return runtime.ToFloat(ctx, min)
	}

	lower, err := floatCeiling(ctx, min)
	if err != nil {
		return runtime.None, err
	}

	upper, err := floatCeiling(ctx, max)
	if err != nil {
		return runtime.None, err
	}

	if lower >= upper {
		return runtime.None, runtime.ArgError(runtime.Error(runtime.ErrInvalidArgument, "interval contains no representable Float"), 1)
	}

	src, err := source(ctx)
	if err != nil {
		return runtime.None, err
	}

	return runtime.Float(interpolateFloat(lower, upper, src.Float64())), nil
}

// Ceiling either endpoint preserves the original half-open interval when an
// Int cannot be represented exactly as Float. Numeric ordering belongs to runtime.
func floatCeiling(ctx context.Context, bound runtime.Value) (float64, error) {
	floating, err := runtime.ToFloat(ctx, bound)
	if err != nil {
		return 0, err
	}

	order, err := runtime.CompareValues(ctx, floating, bound)
	if err != nil {
		return 0, err
	}

	if order == runtime.Less {
		return math.Nextafter(float64(floating), math.Inf(1)), nil
	}

	return float64(floating), nil
}

func interpolateFloat(lower, upper, unit float64) float64 {
	width := upper - lower
	var value float64
	if math.IsInf(width, 0) {
		value = (1-unit)*lower + unit*upper
	} else {
		value = lower + unit*width
	}

	// Rounding can reach the excluded upper endpoint even when unit < 1.
	if value >= upper {
		return math.Nextafter(upper, lower)
	}

	if value < lower {
		return lower
	}

	return value
}
