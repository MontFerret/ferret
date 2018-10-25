package types

import (
	"context"
<<<<<<< HEAD
=======

	"github.com/MontFerret/ferret/pkg/runtime/collections"
>>>>>>> 9f24172... rewrite comments
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

<<<<<<< HEAD
// toArray takes an input value of any type and convert it into an array value.
=======
// ToArray takes an input value of any type and convert it into an array value.
>>>>>>> 9f24172... rewrite comments
// @param (Value) - Input value of arbitrary type.
// @returns (Array)
// None is converted to an empty array
// Boolean values, numbers and strings are converted to an array containing the original value as its single element
// Arrays keep their original value
// Objects / HTML nodes are converted to an array containing their attribute values as array elements.
<<<<<<< HEAD
func ToArray(ctx context.Context, args ...core.Value) (core.Value, error) {
=======
func ToArray(_ context.Context, args ...core.Value) (core.Value, error) {
>>>>>>> 9f24172... rewrite comments
	err := core.ValidateArgs(args, 1, 1)

	if err != nil {
		return values.None, err
	}

	return values.ToArray(ctx, args[0])
}
