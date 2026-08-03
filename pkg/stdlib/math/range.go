package math

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// RANGE returns an array of numbers in the specified range, optionally with increments other than 1.
// @param {Int | Float} start - The value to start the range at (inclusive).
// @param {Int | Float} end - The value to end the range with (inclusive).
// @param {Int | Float} [step=1.0] - How much to increment in every step.
// @return {Int[] | Float[]} - arrayList of numbers in the specified range, optionally with increments other than 1.
func Range(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgs(args, 2, 3); err != nil {
		return runtime.None, err
	}

	if len(args) == 2 {
		return range2(ctx, args[0], args[1])
	}

	return range3(ctx, args[0], args[1], args[2])
}

func range2(ctx context.Context, arg1, arg2 runtime.Value) (runtime.Value, error) {
	return range3(ctx, arg1, arg2, runtime.Float(1))
}

func range3(ctx context.Context, arg1, arg2, arg3 runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgValue(arg1, 0, runtime.AssertNumber); err != nil {
		return runtime.None, err
	}

	if err := runtime.ValidateArgValue(arg2, 1, runtime.AssertNumber); err != nil {
		return runtime.None, err
	}

	if err := runtime.ValidateArgValue(arg3, 2, runtime.AssertNumber); err != nil {
		return runtime.None, err
	}

	step := toFloat(arg3)
	start := toFloat(arg1)
	end := toFloat(arg2)

	arr := runtime.NewArray(int(end))

	for i := start; i <= end; i += step {
		_ = arr.Append(ctx, runtime.NewFloat(i))
	}

	return arr, nil
}
