package types

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

// ToBinary takes an input value of any type and converts it into a binary value.
// @param {Any} value - Input value of arbitrary type.
// @return {Binary} - A binary value.
func ToBinary(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, 1)

	if err != nil {
		return values.None, err
	}

	val := args[0].String()

	return values.NewBinary([]byte(val)), nil
}
