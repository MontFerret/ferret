package objects

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// VALUES return the attribute values of the map as a list.
// @param {Map} map - Target map.
// @return {Any[]} - Values of document returned in any order.
func Values(ctx context.Context, arg runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateType(arg, runtime.TypeMap); err != nil {
		return runtime.None, err
	}

	obj := arg.(runtime.Map)

	return obj.Values(ctx)
}
