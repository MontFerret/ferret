package types

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime"
)

// IS_BINARY checks whether value is a binary value.
// @param {Any} value - Input value of arbitrary type.
// @return {Boolean} - Returns true if value is binary, otherwise false.
func IsBinary(_ context.Context, arg runtime.Value) (runtime.Value, error) {
	return isTypeof(arg, runtime.TypeBinary), nil
}
