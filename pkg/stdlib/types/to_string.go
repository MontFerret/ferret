package types

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

// TO_STRING takes an input value of any type and convert it into a string value.
// @param {Any} value - Input value of arbitrary type.
// @return {String} - String representation of a given value.
func ToString(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, 1)

	if err != nil {
		return values.None, err
	}

	return values.NewString(args[0].String()), nil
}
