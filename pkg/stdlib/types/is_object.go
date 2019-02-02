package types

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

// IsObject checks whether value is an object value.
// @param value (Value) - Input value of arbitrary type.
// @returns (Boolean) - Returns true if value is object, otherwise false.
func IsObject(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, 1)

	if err != nil {
		return values.None, err
	}

	return isTypeof(args[0], types.Object), nil
}
