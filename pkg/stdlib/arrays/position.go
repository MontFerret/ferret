package arrays

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// POSITION returns a value indicating whether an element is contained in array. Optionally returns its position.
// @param {Any[]} array - The source array.
// @param {Any} value - The target value.
// @param {Boolean} [position=False] - Boolean value which indicates whether to return item's position.
// @return {Boolean | Int} - A value indicating whether an element is contained in array.
func Position(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgs(args, 2, 3); err != nil {
		return runtime.None, err
	}

	if len(args) == 2 {
		return position2(ctx, args[0], args[1])
	}

	return position3(ctx, args[0], args[1], args[2])
}

func position2(ctx context.Context, arg1, arg2 runtime.Value) (runtime.Value, error) {
	return position3(ctx, arg1, arg2, runtime.False)
}

func position3(ctx context.Context, arg1, arg2, arg3 runtime.Value) (runtime.Value, error) {
	arr, err := runtime.CastArg[runtime.List](arg1, 0)

	if err != nil {
		return runtime.None, err
	}

	retIdx, err := runtime.CastArg[runtime.Boolean](arg3, 2)

	if err != nil {
		return runtime.None, err
	}

	position, err := arr.IndexOf(ctx, arg2)

	if err != nil {
		return runtime.None, err
	}

	if !retIdx {
		return runtime.NewBoolean(position > -1), nil
	}

	return position, nil
}
