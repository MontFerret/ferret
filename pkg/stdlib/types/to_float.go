package types

import (
	"context"
	"strconv"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

// ToFloat takes an input value of any type and convert it into a float value.
// @param value (Value) - Input value of arbitrary type.
// @returns (Float) -
// None and false are converted to the value 0
// true is converted to 1
// Numbers keep their original value
// Strings are converted to their numeric equivalent if the string contains a valid representation of a number.
// String values that do not contain any valid representation of a number will be converted to the number 0.
// An empty array is converted to 0, an array with one member is converted into the result of TO_NUMBER() for its sole member.
// An array with two or more members is converted to the number 0.
// An object / HTML node is converted to the number 0.
func ToFloat(ctx context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, 1)

	if err != nil {
		return values.None, err
	}

	arg := args[0]

	switch arg.Type() {
	case types.Boolean:
		val := arg.(values.Boolean)

		if val {
			return values.NewFloat(1), nil
		}

		return values.ZeroFloat, nil
	case types.Int:
		val := arg.(values.Int)

		return values.Float(val), nil
	case types.Float:
		return arg, nil
	case types.String:
		str := arg.String()

		if str == "" {
			return values.ZeroFloat, nil
		}

		num, err := strconv.ParseFloat(str, 64)

		if err != nil {
			return values.ZeroFloat, nil
		}

		return values.NewFloat(num), nil
	case types.DateTime:
		val := arg.(values.DateTime)

		if val.IsZero() {
			return values.ZeroFloat, nil
		}

		return values.NewFloat(float64(val.Unix())), nil
	case types.None:
		return values.ZeroFloat, nil
	case types.Array:
		val := arg.(*values.Array)

		if val.Length() == 0 {
			return values.ZeroFloat, nil
		}

		if val.Length() == 1 {
			return ToFloat(ctx, val.Get(0))
		}

		return values.ZeroFloat, nil
	default:
		return values.ZeroFloat, nil
	}
}
