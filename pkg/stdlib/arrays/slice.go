package arrays

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime"
)

// SLICE returns a new sliced array.
// @param {Any[]} array - Source array.
// @param {Int} start - Start position of extraction.
// @param {Int} [length] - Read indicating how many elements to extract.
// @return {Any[]} - Sliced array.
func Slice(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgs(args, 2, 3); err != nil {
		return runtime.None, err
	}

	list, err := runtime.CastList(args[0])

	if err != nil {
		return runtime.None, err
	}

	start, err := runtime.CastInt(args[1])

	if err != nil {
		return runtime.None, err
	}

	var end runtime.Int

	if len(args) > 2 {
		arg3, err := runtime.CastInt(args[2])

		if err != nil {
			return runtime.None, err
		}

		end = start + arg3
	} else {
		size, err := list.Length(ctx)

		if err != nil {
			return runtime.None, err
		}

		end = size
	}

	return list.Slice(ctx, start, end)
}
