package types

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime"
)

// IS_NONE checks whether value is a none value.
// @param {Any} value - Input value of arbitrary type.
// @return {Boolean} - Returns true if value is none, otherwise false.
func IsNone(_ context.Context, arg runtime.Value) (runtime.Value, error) {
	return isTypeof(arg, runtime.TypeNone), nil
}
