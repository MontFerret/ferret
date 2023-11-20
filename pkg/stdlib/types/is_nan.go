package types

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

// IS_NAN checks whether value is NaN.
// @param {Any} value - Input value of arbitrary type.
// @return {Boolean} - Returns true if value is NaN, otherwise false.
func IsNaN(_ context.Context, args ...core.Value) (core.Value, error) {
	if err := core.ValidateArgs(args, 1, 1); err != nil {
		return values.None, err
	}

	val, err := values.CastFloat(args[0])

	if err != nil {
		return values.False, nil
	}

	return values.IsNaN(val), nil
}
