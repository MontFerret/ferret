package types

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// is_list checks whether value is a list value.
// @param value {Any} Input value of arbitrary type.
// @return {Boolean} Returns true if value is list, otherwise false.
func IsList(_ context.Context, arg runtime.Value) (runtime.Value, error) {
	return isTypeof(arg, runtime.TypeList), nil
}
