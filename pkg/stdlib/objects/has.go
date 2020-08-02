package objects

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

// HAS returns the value stored by the given key.
// @param {String} key - The key name string.
// @return {Boolean} - True if the key exists else false.
func Has(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 2, 2)

	if err != nil {
		return values.None, err
	}

	err = core.ValidateType(args[0], types.Object)

	if err != nil {
		return values.None, err
	}

	err = core.ValidateType(args[1], types.String)

	if err != nil {
		return values.None, err
	}

	obj := args[0].(*values.Object)
	keyName := args[1].(values.String)

	_, has := obj.Get(keyName)

	return has, nil
}
