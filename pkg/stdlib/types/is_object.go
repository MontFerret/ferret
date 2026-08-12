package types

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// is_object checks whether value is an object value.
// @param value {Any} Input value of arbitrary type.
// @return {Boolean} Returns true if value is object, otherwise false.
func IsObject(_ context.Context, arg runtime.Value) (runtime.Value, error) {
	return isTypeof(arg, runtime.TypeObject), nil
}
