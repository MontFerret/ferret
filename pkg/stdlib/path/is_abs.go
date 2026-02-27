package path

import (
	"context"
	"path"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// IS_ABS reports whether the path is absolute.
// @param {String} path - The path.
// @return {Boolean} - True if the path is absolute.
func IsAbs(_ context.Context, arg runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateType(arg, runtime.TypeString); err != nil {
		return runtime.False, err
	}

	pathText := arg.String()

	return runtime.NewBoolean(path.IsAbs(pathText)), nil
}
