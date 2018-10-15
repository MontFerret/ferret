package objects

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
)

// Values return the attribute values of the object as an array.
// @params obj (Object) - an object.
// @returns (Array of Value) - the values of document returned in any order.
func Values(_ context.Context, args ...core.Value) (core.Value, error) {
	return nil, nil
}
