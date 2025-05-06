package math

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/internal"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

// SUM returns the sum of the values in a given array.
// @param {Int[] | Float[]} numbers - arrayList of numbers.
// @return {Float} - The sum of the values.
func Sum(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, 1)

	if err != nil {
		return core.None, err
	}

	err = core.ValidateType(args[0], types.Array)

	if err != nil {
		return core.None, err
	}

	arr := args[0].(*internal.Array)

	if arr.Length() == 0 {
		return core.None, nil
	}

	var sum float64

	arr.ForEach(func(value core.Value, idx int) bool {
		err = core.ValidateType(value, types.Int, types.Float)

		if err != nil {
			return false
		}

		sum += toFloat(value)

		return true
	})

	if err != nil {
		return core.None, nil
	}

	return core.NewFloat(sum), nil
}
