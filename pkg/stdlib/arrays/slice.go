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
	_, _, err := runtime.CastVarArgs2[runtime.List, runtime.Int](args)

	if err != nil {
		return runtime.None, err
	}

	if len(args) == 2 {
		return slice2(ctx, args[0], args[1])
	}

	return slice3(ctx, args[0], args[1], args[2])
}

func slice2(ctx context.Context, arg1, arg2 runtime.Value) (runtime.Value, error) {
	return sliceList(ctx, arg1, arg2, runtime.None, false)
}

func slice3(ctx context.Context, arg1, arg2, arg3 runtime.Value) (runtime.Value, error) {
	return sliceList(ctx, arg1, arg2, arg3, true)
}

func sliceList(ctx context.Context, arg1, arg2, arg3 runtime.Value, hasLength bool) (runtime.Value, error) {
	list, err := runtime.CastArg[runtime.List](arg1, 0)
	if err != nil {
		return runtime.None, err
	}

	start, err := runtime.CastArg[runtime.Int](arg2, 1)
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

	if hasLength {
		length, err := runtime.CastArg[runtime.Int](arg3, 2)

		if err != nil {
			return runtime.None, err
		}

		// Handle negative length - return empty array
		if length < 0 {
			return runtime.NewArray(0), nil
		}

		end = start + length
	} else {
		end = size
	}

	// Ensure end doesn't exceed array bounds
	if end > size {
		end = size
	}

	return list.Slice(ctx, start, end)
}
