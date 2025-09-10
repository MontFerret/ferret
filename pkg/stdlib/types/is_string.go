package types

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime"
)

// IS_STRING checks whether value is a string value.
// @param {Any} value - Input value of arbitrary type.
// @return {Boolean} - Returns true if value is string, otherwise false.
func IsString(_ context.Context, arg runtime.Value) (runtime.Value, error) {
	return isTypeof(arg, runtime.TypeString), nil
}
