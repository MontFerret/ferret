package objects

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// Keys returns an independent list of map keys in unspecified order.
// Use sorted(object::keys(value)) when ordering is required.
// @param value {Map} Map whose keys are returned.
// @return {String[]} Independent list of keys in unspecified order.
func Keys(ctx context.Context, arg runtime.Value) (runtime.Value, error) {
	target, err := runtime.CastArg[runtime.Map](arg, 0)
	if err != nil {
		return runtime.None, err
	}

	keys, err := target.Keys(ctx)
	if err != nil {
		return runtime.None, runtime.ArgError(err, 0)
	}

	result, err := copyListValues(ctx, keys)
	if err != nil {
		return runtime.None, runtime.ArgError(err, 0)
	}

	return result, nil
}
