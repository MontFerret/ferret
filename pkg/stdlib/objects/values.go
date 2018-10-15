package objects

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/values"

	"github.com/MontFerret/ferret/pkg/runtime/core"
)

// Values return the attribute values of the object as an array.
// @params obj (Object) - an object.
// @returns (Array of Value) - the values of document returned in any order.
func Values(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, 1)
	if err != nil {
		return values.None, err
	}

	err = core.ValidateType(args[0], core.ObjectType)
	if err != nil {
		return values.None, err
	}

	obj := args[0].(*values.Object)
	vals := values.NewArray(0)

	obj.ForEach(func(val core.Value, key string) bool {
		if values.IsCloneable(val) {
			val = val.(core.Cloneable).Clone()
		}
		vals.Push(val)
		return true
	})

	return vals, nil
}
