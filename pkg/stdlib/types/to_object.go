package types

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// to_object converts the given value to an object.
// @param value {Any} Input value of arbitrary type.
// @return {Map} Returns the object representation of the given value.
func ToObject(ctx context.Context, arg runtime.Value) (runtime.Value, error) {
	return runtime.ToMap(ctx, arg)
}
