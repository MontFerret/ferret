package types

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime"
)

// IS_FLOAT checks whether value is a float value.
// @param {Any} value - Input value of arbitrary type.
// @return {Boolean} - Returns true if value is float, otherwise false.
func IsFloat(_ context.Context, arg runtime.Value) (runtime.Value, error) {
	return isTypeof(arg, runtime.TypeFloat), nil
}
