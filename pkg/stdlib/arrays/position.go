package arrays

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

// POSITION returns a value indicating whether an element is contained in array. Optionally returns its position.
// @param {Any[]} array - The source array.
// @param {Any} value - The target value.
// @param {Boolean} [position=False] - Boolean value which indicates whether to return item's position.
// @return {Boolean | Int} - A value indicating whether an element is contained in array.
func Position(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 2, 3)

	if err != nil {
		return values.None, err
	}

	err = values.AssertArray(args[0])

	if err != nil {
		return values.None, err
	}

	arr := args[0].(*values.Array)
	el := args[1]
	retIdx := false

	if len(args) > 2 {
		err = values.AssertBoolean(args[2])

		if err != nil {
			return values.None, err
		}

		retIdx = values.Compare(args[2], values.True) == 0
	}

	position := arr.IndexOf(el)

	if !retIdx {
		return values.NewBoolean(position > -1), nil
	}

	return position, nil
}
