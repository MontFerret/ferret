package types

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

// ToDateTime takes an input value of any type and converts it into the appropriate date time value.
// @param value (Value) - Input value of arbitrary type.
// @returns (DateTime) - Parsed date time.
func ToDateTime(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, 1)

	if err != nil {
		return values.None, err
	}

	return values.ParseDateTime(args[0].String())
}
