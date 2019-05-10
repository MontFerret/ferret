package types

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
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

	return values.ToBoolean(args[0]), nil
}
