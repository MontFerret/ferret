package types

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime"
)

// IS_LSIT checks whether value is a list value.
// @param {Any} value - Input value of arbitrary type.
// @return {Boolean} - Returns true if value is list, otherwise false.
func IsList(_ context.Context, arg runtime.Value) (runtime.Value, error) {
	return isTypeof(arg, runtime.TypeList), nil
}
