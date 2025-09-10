package types

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime"
)

// TO_BOOL takes an input value of any type and converts it into the appropriate boolean value.
// None is converted to false
// Numbers are converted to true, except for 0, which is converted to false
// Strings are converted to true if they are non-empty, and to false otherwise
// Dates are converted to true if they are not zero, and to false otherwise
// Arrays are always converted to true (even if empty)
// Objects / HtmlNodes / Binary are always converted to true
// @param {Any} value - Input value of arbitrary type.
// @return {Boolean} - The appropriate boolean value.
func ToBool(_ context.Context, arg runtime.Value) (runtime.Value, error) {
	return runtime.ToBoolean(arg), nil
}
