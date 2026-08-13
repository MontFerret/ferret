package types

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// is_string checks whether value is a string value.
// @param value {Any} Input value of arbitrary type.
// @return {Boolean} Returns true if value is string, otherwise false.
func IsString(_ context.Context, arg runtime.Value) (runtime.Value, error) {
	return isTypeof(arg, runtime.TypeString), nil
}
