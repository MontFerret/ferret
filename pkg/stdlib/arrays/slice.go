package arrays

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

// SLICE returns a new sliced array.
// @param {Any[]} array - Source array.
// @param {Int} start - Start position of extraction.
// @param {Int} [length] - Read indicating how many elements to extract.
// @return {Any[]} - Sliced array.
func Slice(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 2, 3)

	if err != nil {
		return values.None, err
	}

	err = values.AssertArray(args[0])

	if err != nil {
		return values.None, err
	}

	err = values.AssertInt(args[1])

	if err != nil {
		return values.None, err
	}

	arr := args[0].(*values.Array)
	start := int(args[1].(values.Int))
	length := arr.Length()

	if len(args) > 2 {
		lengthArg, ok := args[2].(values.Int)

		if ok && lengthArg > 0 {
			length = start + int(lengthArg)
		}
	}

	return arr.Slice(start, length), nil
}
