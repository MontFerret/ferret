package types

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime"
)

// IS_NAN checks whether value is NaN.
// @param {Any} value - Input value of arbitrary type.
// @return {Boolean} - Returns true if value is NaN, otherwise false.
func IsNaN(_ context.Context, args ...runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgs(args, 1, 1); err != nil {
		return runtime.None, err
	}

	val, err := runtime.CastFloat(args[0])

	if err != nil {
		return runtime.False, nil
	}

	return runtime.IsNaN(val), nil
}
