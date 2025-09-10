package types

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime"
)

// IS_BOOL checks whether value is a boolean value.
// @param {Any} value - Input value of arbitrary type.
// @return {Boolean} - Returns true if value is boolean, otherwise false.
func IsBool(_ context.Context, arg runtime.Value) (runtime.Value, error) {
	return isTypeof(arg, runtime.TypeBoolean), nil
}
