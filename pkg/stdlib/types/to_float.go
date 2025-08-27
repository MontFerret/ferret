package types

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime"
)

// TO_FLOAT takes an input value of any type and convert it into a float value.
// None and false are converted to the value 0
// true is converted to 1
// Numbers keep their original value
// Strings are converted to their numeric equivalent if the string contains a valid representation of a number.
// String values that do not contain any valid representation of a number will be converted to the number 0.
// An empty array is converted to 0, an array with one member is converted into the result of TO_NUMBER() for its sole member.
// An array with two or more members is converted to the number 0.
// An object / HTML node is converted to the number 0.
// @param {Any} value - Input value of arbitrary type.
// @return {Float} - A float value.
func ToFloat(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
	err := runtime.ValidateArgs(args, 1, 1)

	if err != nil {
		return runtime.None, err
	}

	return runtime.ToFloat(ctx, args[0])
}
