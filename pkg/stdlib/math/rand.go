package math

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// rand return a pseudo-random number between 0 and 1.
// @param max {Int | Float} Upper limit.
// @param min {Int | Float} Lower limit.
// @return {Float} A number greater than 0 and less than 1.
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

// rand return a pseudo-random number between 0 and 1.
// @return {Float} A number greater than 0 and less than 1.
func rand0(context.Context) (runtime.Value, error) {
	return runtime.NewFloat(runtime.RandomDefault()), nil
}

// rand return a pseudo-random number between 0 and 1.
// @param max {Int | Float} Upper limit.
// @return {Float} A number greater than 0 and less than 1.
func rand1(ctx context.Context, arg1 runtime.Value) (runtime.Value, error) {
	max, err := runtime.ToFloat(ctx, arg1)
	if err != nil {
		return runtime.None, err
	}

	upper, lower := runtime.NumberBoundaries(float64(max))

	return runtime.NewFloat(runtime.Random(upper, lower)), nil
}

// rand return a pseudo-random number between 0 and 1.
// @param max {Int | Float} Upper limit.
// @param min {Int | Float} Lower limit.
// @return {Float} A number greater than 0 and less than 1.
func rand2(ctx context.Context, arg1, arg2 runtime.Value) (runtime.Value, error) {
	max, err := runtime.ToFloat(ctx, arg1)

	if err != nil {
		return runtime.None, err
	}

	min, err := runtime.ToFloat(ctx, arg2)
	if err != nil {
		return runtime.None, err
	}

	return runtime.NewFloat(runtime.Random(float64(max), float64(min))), nil
}
