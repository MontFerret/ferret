package types

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime"
)

// IS_INT checks whether value is a int value.
// @param {Any} value - Input value of arbitrary type.
// @return {Boolean} - Returns true if value is int, otherwise false.
func IsInt(_ context.Context, arg runtime.Value) (runtime.Value, error) {
	return isTypeof(arg, runtime.TypeInt), nil
}
