package objects

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// HAS returns the value stored by the given key.
// @param {String} key - The key name string.
// @return {Boolean} - True if the key exists else false.
func Has(ctx context.Context, arg1, arg2 runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgType(arg1, 0, runtime.TypeMap); err != nil {
		return runtime.None, err
	}

	if err := runtime.ValidateArgType(arg2, 1, runtime.TypeString); err != nil {
		return runtime.None, err
	}

	target := arg1.(runtime.Map)
	key := arg2.(runtime.String)

	return target.ContainsKey(ctx, key)
}
