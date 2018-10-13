package math

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/pkg/errors"
	"math"
)

/*
 * Return the nth percentile of the values in a given array.
 * @param array (Array) - Array of numbers.
 * @param numb (Int) - A number which must be between 0 (excluded) and 100 (included).
 * @param method (String, optional) - "rank" (default) or "interpolation".
 * @returns (Float) - The nth percentile, or null if the array is empty or only null values are contained in it or the percentile cannot be calculated.
 */
func Percentile(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 2, 3)

	if err != nil {
		return values.None, err
	}

	err = core.ValidateType(args[0], core.ArrayType)

	if err != nil {
		return values.None, err
	}

	err = core.ValidateType(args[1], core.IntType)

	if err != nil {
		return values.None, err
	}

	// TODO: Implement different methods
	//method := "rank"
	//
	//if len(args) > 2 {
	//	err = core.ValidateType(args[2], core.StringType)
	//
	//	if err != nil {
	//		return values.None, err
	//	}
	//
	//	if args[2].String() == "interpolation" {
	//		method = "interpolation"
	//	}
	//}

	arr := args[0].(*values.Array)
	percent := values.Float(args[1].(values.Int))

	if arr.Length() == 0 {
		return values.NewFloat(math.NaN()), nil
	}

	if percent <= 0 || percent > 100 {
		return values.NewFloat(math.NaN()), errors.New("input is outside of range")
	}

	sorted := arr.Sort()

	// Multiply percent by length of input
	l := values.Float(sorted.Length())
	index := (percent / 100) * l
	even := values.Float(values.Int(index))

	var percentile core.Value

	// Check if the index is a whole number
	if index == even {

		// Convert float to int
		i := values.Int(index)

		// Find the value at the index
		percentile = sorted.Get(i - 1)
	} else if index > 1 {

		// Convert float to int via truncation
		i := values.Int(index)

		// Find the average of the index and following values
		percentile, _ = mean(values.NewArrayWith(sorted.Get(i-1), sorted.Get(i)))

	} else {
		return values.NewFloat(math.NaN()), errors.New("input is outside of range")
	}

	return percentile, nil
}
