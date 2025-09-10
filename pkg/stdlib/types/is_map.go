package types

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime"
)

// IS_MAP checks whether value is a map value.
// @param {Any} value - Input value of arbitrary type.
// @return {Boolean} - Returns true if value is map, otherwise false.
func IsMap(_ context.Context, arg runtime.Value) (runtime.Value, error) {
	return isTypeof(arg, runtime.TypeMap), nil
}
