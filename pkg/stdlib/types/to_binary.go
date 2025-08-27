package types

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime"
)

// ToBinary takes an input value of any type and converts it into a binary value.
// @param {Any} value - Input value of arbitrary type.
// @return {Binary} - A binary value.
func ToBinary(_ context.Context, args ...runtime.Value) (runtime.Value, error) {
	err := runtime.ValidateArgs(args, 1, 1)

	if err != nil {
		return runtime.None, err
	}

	val := args[0].String()

	return runtime.NewBinary([]byte(val)), nil
}
