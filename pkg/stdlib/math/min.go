package math

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

// Min returns the smallest (arithmetic mean) of the values in array.
// @param array (Array) - Array of numbers.
// @returns (Float) - The smallest of the values in array.
func Min(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, 1)

	if err != nil {
		return values.None, err
	}

	err = core.ValidateType(args[0], types.Array)

	if err != nil {
		return values.None, err
	}

	arr := args[0].(*values.Array)

	if arr.Length() == 0 {
		return values.None, nil
	}

	var min float64

	arr.ForEach(func(value core.Value, idx int) bool {
		err = core.ValidateType(value, types.Int, types.Float)

		if err != nil {
			return false
		}

		fv := toFloat(value)

		if min > fv || idx == 0 {
			min = fv
		}

		return true
	})

	if err != nil {
		return values.None, nil
	}

	return values.NewFloat(min), nil
}
