package arrays

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

// Slice returns a new sliced array.
// @param array (Array) - Source array.
// @param start (Int) - Start position of extraction.
// @param length (Int, optional) - Read indicating how many elements to extract.
// @returns (Array) - Sliced array.
func Slice(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 2, 3)

	if err != nil {
		return values.None, err
	}

	err = core.ValidateType(args[0], types.Array)

	if err != nil {
		return values.None, err
	}

	err = core.ValidateType(args[1], types.Int)

	if err != nil {
		return values.None, err
	}

	arr := args[0].(*values.Array)
	start := args[1].(values.Int)
	length := values.NewInt(int(arr.Length()))

	if len(args) > 2 {
		if args[2].Type() == types.Int {
			arg2 := args[2].(values.Int)

			if arg2 > 0 {
				length = start + args[2].(values.Int)
			}
		}
	}

	return arr.Slice(start, length), nil
}
