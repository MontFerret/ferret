package types

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

// ToBool takes an input value of any type and converts it into the appropriate boolean value.
// @param value (Value) - Input value of arbitrary type.
// @returns (Boolean) -
// None is converted to false
// Numbers are converted to true, except for 0, which is converted to false
// Strings are converted to true if they are non-empty, and to false otherwise
// Dates are converted to true if they are not zero, and to false otherwise
// Arrays are always converted to true (even if empty)
// Objects / HtmlNodes / Binary are always converted to true
func ToBool(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, 1)

	if err != nil {
		return values.None, err
	}

	arg := args[0]

	switch arg.Type() {
	case types.Boolean:
		return arg, nil
	case types.Int:
		val := arg.(values.Int)

		if val != 0 {
			return values.True, nil
		}

		return values.False, nil
	case types.Float:
		val := arg.(values.Float)

		if val != 0 {
			return values.True, nil
		}

		return values.False, nil
	case types.String:
		if arg.String() != "" {
			return values.True, nil
		}

		return values.False, nil
	case types.DateTime:
		val := arg.(values.DateTime)

		if !val.IsZero() {
			return values.True, nil
		}

		return values.False, nil
	case types.None:
		return values.False, nil
	default:
		return values.True, nil
	}
}
