package types

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

// ToBinary takes an input value of any type and convert it into a binary value.
// @param value (Value) - Input value of arbitrary type.
// @return (Binary) - String representation of a given value.
func ToBinary(ctx context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, 1)

	if err != nil {
		return values.None, err
	}

	val := args[0].String()

	return values.NewBinary([]byte(val)), nil
}
