package types

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime"
)

// IS_ARRAY checks whether value is an array value.
// @param {Any} value - Input value of arbitrary type.
// @return {Boolean} - Returns true if value is array, otherwise false.
func IsArray(_ context.Context, arg runtime.Value) (runtime.Value, error) {
	return isTypeof(arg, runtime.TypeArray), nil
}
