package objects

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

// VALUES return the attribute values of the object as an array.
// @param {Object} object - Target object.
// @return {Any[]} - Values of document returned in any order.
func Values(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, 1)

	if err != nil {
		return values.None, err
	}

	err = core.ValidateType(args[0], types.Object)

	if err != nil {
		return values.None, err
	}

	obj := args[0].(*values.Object)
	vals := values.NewArray(0)

	obj.ForEach(func(val core.Value, key string) bool {
		cloneable, ok := val.(core.Cloneable)

		if ok {
			val = cloneable.Clone()
		}

		vals.Push(val)

		return true
	})

	return vals, nil
}
