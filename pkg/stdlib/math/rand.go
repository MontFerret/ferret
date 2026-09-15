package math

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/rnd"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// Rand preserves the legacy zero-, one-, and two-argument random calculations
// using the Session source transported by ctx.
// Deprecated: Use random.Float, random.Int, or random.Bool with an execution context.
func Rand(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgs(args, 0, 2); err != nil {
		return runtime.None, err
	}

	switch len(args) {
	case 0:
		return rand0(ctx)
	case 1:
		return rand1(ctx, args[0])
	default:
		return rand2(ctx, args[0], args[1])
	}
}

// rand draws a non-cryptographic pseudo-random Float from the Session source.
// @return {Float} A number in [0, 1).
// @deprecated Use random::float() instead.
func rand0(ctx context.Context) (runtime.Value, error) {
	src, ok := rnd.FromContext(ctx)
	if !ok {
		return runtime.None, runtime.Error(runtime.ErrUnexpected, "random source missing from execution context")
	}

	return runtime.Float(src.Float64()), nil
}

// rand applies the historical rounded random calculation around a supplied value.
// @param max {Int | Float} Converted to Float; the historical lower bound is max/2 and upper bound is max*2.
// @return {Float} The historical result floor(u*(max*2-max/2+1))+max/2, where u is in [0, 1).
// @deprecated Use random::int(min, max) for inclusive integer bounds or random::float(min, max) for continuous bounds. Supply explicit bounds; the legacy one-argument form uses max/2 and max*2.
func rand1(ctx context.Context, arg1 runtime.Value) (runtime.Value, error) {
	max, err := runtime.ToFloat(ctx, arg1)
	if err != nil {
		return runtime.None, err
	}

	upper, lower := runtime.NumberBoundaries(float64(max))

	return randRange(ctx, upper, lower)
}

// rand applies the historical rounded random calculation with maximum first.
// @param max {Int | Float} Upper bound, converted to Float without canonical range validation.
// @param min {Int | Float} Lower bound, converted to Float without canonical range validation.
// @return {Float} The historical result floor(u*(max-min+1))+min, where u is in [0, 1).
// @deprecated Use random::int(min, max) for inclusive integer bounds or random::float(min, max) for continuous bounds. Reverse the legacy argument order; canonical functions validate bounds and types.
func rand2(ctx context.Context, arg1, arg2 runtime.Value) (runtime.Value, error) {
	max, err := runtime.ToFloat(ctx, arg1)
	if err != nil {
		return runtime.None, err
	}

	min, err := runtime.ToFloat(ctx, arg2)
	if err != nil {
		return runtime.None, err
	}

	return randRange(ctx, float64(max), float64(min))
}

func randRange(ctx context.Context, max, min float64) (runtime.Value, error) {
	src, ok := rnd.FromContext(ctx)
	if !ok {
		return runtime.None, runtime.Error(runtime.ErrUnexpected, "random source missing from execution context")
	}

	return runtime.Float(src.LegacyFloat64(max, min)), nil
}
