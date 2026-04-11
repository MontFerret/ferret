package path

import (
	"context"
	"path"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// CLEAN returns the shortest path name equivalent to path.
// @param {String} path - The path.
// @return {String} - The shortest path name equivalent to path
func Clean(_ context.Context, arg runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgType(arg, 0, runtime.TypeString); err != nil {
		return runtime.EmptyString, err
	}

	pathText := arg.String()

	return runtime.NewString(path.Clean(pathText)), nil
}
