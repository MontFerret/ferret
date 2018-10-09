package objects

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

/*
 * Merge the given objects into a single object.
 * @params objs (Array Of Object OR Objects) - objects to merge.
 * @returns (Object) - Object created by merging.
 */
func Merge(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, core.MaxArgs)

	if err != nil {
		return values.None, err
	}

	objs := values.NewArrayWith(args...)

	if len(args) == 1 && args[0].Type() == core.ArrayType {
		objs = args[0].(*values.Array)
	}

	err = validateTopLevelObjects(objs)

	if err != nil {
		return values.None, err
	}

	return mergeArray(objs), nil
}

func validateTopLevelObjects(arr *values.Array) error {
	return nil
}

func mergeArray(arr *values.Array) *values.Object {
	return values.NewObject()
}
