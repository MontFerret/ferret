package path

import (
	"context"
	"path"

	"github.com/MontFerret/ferret/pkg/runtime"
)

// IS_ABS reports whether the path is absolute.
// @param {String} path - The path.
// @return {Boolean} - True if the path is absolute.
func IsAbs(_ context.Context, args ...runtime.Value) (runtime.Value, error) {
	err := runtime.ValidateArgs(args, 1, 1)

	if err != nil {
		return runtime.False, err
	}

	err = runtime.ValidateType(args[0], runtime.TypeString)

	if err != nil {
		return runtime.False, err
	}

	pathText := args[0].String()

	return runtime.NewBoolean(path.IsAbs(pathText)), nil
}
