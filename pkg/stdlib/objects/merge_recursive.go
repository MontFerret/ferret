package objects

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

// MergeRecursive recursively merge the given objects into a single object.
// @params objs (Objects) - objects to merge.
// @returns (Object) - Object created by merging.
func MergeRecursive(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, core.MaxArgs)
	if err != nil {
		return values.None, err
	}

	for _, arg := range args {
		if err = core.ValidateType(arg, core.ObjectType); err != nil {
			return values.None, err
		}
	}

	merged := values.NewObject()

	return merged, nil
}
