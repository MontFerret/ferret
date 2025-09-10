package types

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime"
)

// IS_NAN checks whether value is NaN.
// @param {Any} value - Input value of arbitrary type.
// @return {Boolean} - Returns true if value is NaN, otherwise false.
func IsNaN(_ context.Context, arg runtime.Value) (runtime.Value, error) {
	val, err := runtime.CastFloat(arg)

	if err != nil {
		return runtime.False, nil
	}

	return runtime.IsNaN(val), nil
}
