package types

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// to_number takes an input value of any type and convert it into a number value.
// None and false are converted to the zero float
// true is converted to 1 float
// Numbers keep their original value
// Strings are converted to their numeric equivalent if the string contains a valid representation of a number.
// String values that do not contain any valid representation of a number will be converted to the number 0.
// An empty array is converted to 0, an array with one member is converted into the result of to_number() for its sole member.
// An array with two or more members is converted to the number 0.
// An object / HTML node is converted to the number 0.
// @param value {Any} Input value of arbitrary type.
// @return {Int|Float} An integer or float value.
func ToNumber(ctx context.Context, arg runtime.Value) (runtime.Value, error) {
	return runtime.ToNumber(ctx, arg)
}
