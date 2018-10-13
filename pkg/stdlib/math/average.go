package math

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

/*
 * Returns the average (arithmetic mean) of the values in array.
 * @param array (Array) - Array of numbers.
 * @returns (Float) - The average (arithmetic mean) of the values in array.
 */
func Average(_ context.Context, args ...core.Value) (core.Value, error) {
	var err error
	err = core.ValidateArgs(args, 1, 1)

	if err != nil {
		return values.None, err
	}

	err = core.ValidateType(args[0], core.ArrayType)

	if err != nil {
		return values.None, err
	}

	arr := args[0].(*values.Array)

	var sum float64

	arr.ForEach(func(value core.Value, idx int) bool {
		err = core.ValidateType(value, core.FloatType, core.IntType)

		if err != nil {
			return false
		}

		var fv float64

		if value.Type() == core.IntType {
			fv = float64(value.(values.Int))
		} else {
			fv = float64(value.(values.Float))
		}

		sum += fv

		return true
	})

	if err != nil {
		return values.None, err
	}

	count := arr.Length()

	return values.Float(sum / float64(count)), nil
}
