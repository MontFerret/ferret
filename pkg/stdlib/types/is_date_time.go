package types

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime"
)

// IS_DATETIME checks whether value is a date time value.
// @param {Any} value - Input value of arbitrary type.
// @return {Boolean} - Returns true if value is date time, otherwise false.
func IsDateTime(_ context.Context, arg runtime.Value) (runtime.Value, error) {
	return isTypeof(arg, runtime.TypeDateTime), nil
}
