package types

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

// IsNaN checks whether value is NaN.
// @param value (Value) - Input value of arbitrary type.
// @returns (Boolean) - Returns true if value is NaN, otherwise false.
func IsNaN(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, 1)

	if err != nil {
		return values.None, err
	}

	if args[0].Type() != types.Float {
		return values.False, nil
	}

	return values.IsNaN(args[0].(values.Float)), nil
}
