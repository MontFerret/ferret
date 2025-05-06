package arrays

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/internal"

	"github.com/MontFerret/ferret/pkg/runtime/core"
)

// POSITION returns a value indicating whether an element is contained in array. Optionally returns its position.
// @param {Any[]} array - The source array.
// @param {Any} value - The target value.
// @param {Boolean} [position=False] - Boolean value which indicates whether to return item's position.
// @return {Boolean | Int} - A value indicating whether an element is contained in array.
func Position(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 2, 3)

	if err != nil {
		return core.None, err
	}

	err = core.AssertList(args[0])

	if err != nil {
		return core.None, err
	}

	arr := args[0].(*internal.Array)
	el := args[1]
	retIdx := false

	if len(args) > 2 {
		err = core.AssertBoolean(args[2])

		if err != nil {
			return core.None, err
		}

		retIdx = core.CompareValues(args[2], core.True) == 0
	}

	position := arr.IndexOf(el)

	if !retIdx {
		return core.NewBoolean(position > -1), nil
	}

	return position, nil
}
