package types

import (
	"context"
	"strconv"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

// ToInt takes an input value of any type and convert it into an integer value.
// @param value (Value) - Input value of arbitrary type.
// @returns (Int) -
// None and false are converted to the value 0
// true is converted to 1
// Numbers keep their original value
// Strings are converted to their numeric equivalent if the string contains a valid representation of a number.
// String values that do not contain any valid representation of a number will be converted to the number 0.
// An empty array is converted to 0, an array with one member is converted into the result of TO_NUMBER() for its sole member.
// An array with two or more members is converted to the number 0.
// An object / HTML node is converted to the number 0.
func ToInt(ctx context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, 1)

	if err != nil {
		return values.None, err
	}

	arg := args[0]

	switch arg.Type() {
	case core.BooleanType:
		val := arg.(values.Boolean)

		if val {
			return values.NewInt(1), nil
		}

		return values.ZeroInt, nil
	case core.IntType:
		return arg, nil
	case core.FloatType:
		val := arg.(values.Float)

		return values.Int(val), nil
	case core.StringType:
		str := arg.String()

		if str == "" {
			return values.ZeroInt, nil
		}

		num, err := strconv.Atoi(str)

		if err != nil {
			return values.ZeroInt, nil
		}

		return values.NewInt(int64(num)), nil
	case core.DateTimeType:
		val := arg.(values.DateTime)

		if val.IsZero() {
			return values.ZeroInt, nil
		}

		return values.NewInt(val.Unix()), nil
	case core.NoneType:
		return values.ZeroInt, nil
	case core.ArrayType:
		val := arg.(*values.Array)

		if val.Length() == 0 {
			return values.ZeroInt, nil
		}

		if val.Length() == 1 {
			return ToInt(ctx, val.Get(0))
		}

		return values.ZeroInt, nil
	default:
		return values.ZeroInt, nil
	}
}
