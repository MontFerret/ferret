package arrays

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

/*
 * Returns a value indicating whether an element is contained in array. Optionally returns its position.
 * @param array (Array) - The source array.
 * @param value (Value) - The target value.
 * @param returnIndex (Boolean, optional) - Value which indicates whether to return item's position.
 */
func Position(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 2, 3)

	if err != nil {
		return values.None, err
	}

	err = core.ValidateType(args[0], core.ArrayType)

	if err != nil {
		return values.None, err
	}

	arr := args[0].(*values.Array)
	el := args[1]
	retIdx := false

	if len(args) > 2 {
		err = core.ValidateType(args[2], core.BooleanType)

		if err != nil {
			return values.None, err
		}

		retIdx = args[2].Compare(values.True) == 0
	}

	position := arr.IndexOf(el)

	if !retIdx {
		return values.NewBoolean(position > -1), nil
	}

	return position, nil
}
