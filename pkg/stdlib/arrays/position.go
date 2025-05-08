package arrays

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime"
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

	arr, err := runtime.CastList(args[0])

	if err != nil {
		return runtime.None, err
	}

	el := args[1]
	var retIdx runtime.Boolean

	if len(args) > 2 {
		arg2, err := runtime.CastBoolean(args[2])

		if err != nil {
			return runtime.None, err
		}

		retIdx = arg2
	}

	position, err := arr.IndexOf(ctx, el)

	if err != nil {
		return runtime.None, err
	}

	if !retIdx {
		return runtime.NewBoolean(position > -1), nil
	}

	return position, nil
}
