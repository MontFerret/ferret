package arrays

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// SLICE returns a new sliced array.
// @param {Any[]} array - Source array.
// @param {Int} start - Start position of extraction.
// @param {Int} [length] - Read indicating how many elements to extract.
// @return {Any[]} - Sliced array.
func Slice(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
	list, start, err := runtime.CastVarArgs2[runtime.List, runtime.Int](args)

	if err != nil {
		return runtime.None, err
	}

	size, err := list.Length(ctx)

	if err != nil {
		return runtime.None, err
	}

	// Handle negative start index - return empty array
	if start < 0 {
		return runtime.NewArray(0), nil
	}

	// Handle start index beyond array length - return empty array
	if start >= size {
		return runtime.NewArray(0), nil
	}

	var end runtime.Int

	if len(args) > 2 {
		arg3, err := runtime.CastArgAt[runtime.Int](args, 2)

		if err != nil {
			return runtime.None, err
		}

		// Handle negative length - return empty array
		if arg3 < 0 {
			return runtime.NewArray(0), nil
		}

		end = start + arg3
	} else {
		end = size
	}

	// Ensure end doesn't exceed array bounds
	if end > size {
		end = size
	}

	return list.Slice(ctx, start, end)
}
