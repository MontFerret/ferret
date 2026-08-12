package types

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// is_array checks whether value is an array value.
// @param value {Any} Input value of arbitrary type.
// @return {Boolean} Returns true if value is array, otherwise false.
func IsArray(_ context.Context, arg runtime.Value) (runtime.Value, error) {
	return isTypeof(arg, runtime.TypeArray), nil
}
