package types

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime"
)

// TO_OBJECT converts the given value to an object.
// @param {Any} value - Input value of arbitrary type.
// @return {Object} - Returns the object representation of the given value.
func ToObject(ctx context.Context, arg runtime.Value) (runtime.Value, error) {
	return runtime.ToMap(ctx, arg), nil
}
