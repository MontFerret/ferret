package types

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// ToBinary takes an input value of any type and converts it into a binary value.
// @param {Any} value - Input value of arbitrary type.
// @return {binary} - A binary value.
func ToBinary(_ context.Context, arg runtime.Value) (runtime.Value, error) {
	val := arg.String()

	return runtime.NewBinary([]byte(val)), nil
}
