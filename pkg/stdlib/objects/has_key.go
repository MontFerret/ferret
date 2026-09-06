package objects

import (
	"context"
	"fmt"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// HasKey checks key existence, including keys whose value is none.
// @param object {Map} Object or map to inspect.
// @param key {String} The key name string.
// @return {Boolean} True if the key exists else false.
func HasKey(ctx context.Context, arg1, arg2 runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgType(arg1, 0, runtime.TypeMap); err != nil {
		return runtime.None, err
	}

	if err := runtime.ValidateArgType(arg2, 1, runtime.TypeString); err != nil {
		return runtime.None, err
	}

	target := arg1.(runtime.Map)
	key := arg2.(runtime.String)

	if err := ctx.Err(); err != nil {
		return runtime.None, err
	}

	exists, err := target.ContainsKey(ctx, key)
	if err != nil {
		return runtime.None, runtime.ArgError(fmt.Errorf("key %q: %w", key, err), 0)
	}

	return exists, nil
}
