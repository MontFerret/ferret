package path

import (
	"context"
	"path"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// DIR returns the directory component of path.
// @param {String} path - The path.
// @return {String} - The directory component of path.
func Dir(_ context.Context, arg runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgType(arg, 0, runtime.TypeString); err != nil {
		return runtime.EmptyString, err
	}

	pathText := arg.String()

	return runtime.NewString(path.Dir(pathText)), nil
}
