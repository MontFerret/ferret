package types

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

// TypeName returns the data type name of value.
// @param value (Value) - Input value of arbitrary type.
// @returns (Boolean) - Returns string representation of a type.
func TypeName(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, 1)

	if err != nil {
		return values.None, err
	}

	return values.NewString(args[0].Type().String()), nil
}
