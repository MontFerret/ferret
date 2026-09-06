package objects

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// Values returns cloned or copied map values in unspecified order.
// @param value {Map} Map whose values are returned.
// @return {Any[]} Independent list of values, following each value's clone or copy contract.
func Values(ctx context.Context, arg runtime.Value) (runtime.Value, error) {
	target, err := runtime.CastArg[runtime.Map](arg, 0)
	if err != nil {
		return runtime.None, err
	}

	values, err := target.Values(ctx)
	if err != nil {
		return runtime.None, runtime.ArgError(err, 0)
	}

	result, err := copyListValues(ctx, values)
	if err != nil {
		return runtime.None, runtime.ArgError(err, 0)
	}

	return result, nil
}
