package path

import (
	"context"
	"path"

	"github.com/MontFerret/ferret/pkg/runtime"
)

// CLEAN returns the shortest path name equivalent to path.
// @param {String} path - The path.
// @return {String} - The shortest path name equivalent to path
func Clean(_ context.Context, args ...runtime.Value) (runtime.Value, error) {
	err := runtime.ValidateArgs(args, 1, 1)

	if err != nil {
		return runtime.EmptyString, err
	}

	err = runtime.ValidateType(args[0], runtime.TypeString)

	if err != nil {
		return runtime.None, err
	}

	pathText := args[0].String()

	return runtime.NewString(path.Clean(pathText)), nil
}
