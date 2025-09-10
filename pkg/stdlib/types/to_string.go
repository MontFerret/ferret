package types

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime"
)

// TO_STRING takes an input value of any type and convert it into a string value.
// @param {Any} value - Input value of arbitrary type.
// @return {String} - String representation of a given value.
func ToString(_ context.Context, arg runtime.Value) (runtime.Value, error) {
	return runtime.NewString(arg.String()), nil
}
