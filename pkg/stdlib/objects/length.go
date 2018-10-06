package objects

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

/*
 * Returns the Int equal to the number of values inside the object.
 * @params (Object) - The source object.
 * @returns (Int) - New Int equal to the number of values inside the source object.
 */
func Length(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, 1)

	if err != nil {
		return values.None, err
	}

	err = core.ValidateType(args[0], core.ObjectType)

	if err != nil {
		return values.None, err
	}

	obj := args[0].(*values.Object)

	return values.NewInt(int(obj.Length())), nil
}
