package objects

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime"
)

// HAS returns the value stored by the given key.
// @param {String} key - The key name string.
// @return {Boolean} - True if the key exists else false.
// TODO: REWRITE TO USE LIST & MAP instead
func Has(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
	err := runtime.ValidateArgs(args, 2, 2)

	if err != nil {
		return runtime.None, err
	}

	err = runtime.ValidateType(args[0], runtime.TypeObject)

	if err != nil {
		return runtime.None, err
	}

	err = runtime.ValidateType(args[1], runtime.TypeString)

	if err != nil {
		return runtime.None, err
	}

	obj := args[0].(*runtime.Object)
	keyName := args[1].(runtime.String)

	return obj.ContainsKey(ctx, keyName)
}
