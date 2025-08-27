package types

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime"
)

// TO_STRING takes an input value of any type and convert it into a string value.
// @param {Any} value - Input value of arbitrary type.
// @return {String} - String representation of a given value.
func ToString(_ context.Context, args ...runtime.Value) (runtime.Value, error) {
	err := runtime.ValidateArgs(args, 1, 1)

	if err != nil {
		return runtime.None, err
	}

	return runtime.NewString(args[0].String()), nil
}
