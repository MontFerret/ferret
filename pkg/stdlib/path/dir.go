package path

import (
	"context"
	"path"

	"github.com/MontFerret/ferret/pkg/runtime"
)

// DIR returns the directory component of path.
// @param {String} path - The path.
// @return {String} - The directory component of path.
func Dir(_ context.Context, args ...runtime.Value) (runtime.Value, error) {
	err := runtime.ValidateArgs(args, 1, 1)

	if err != nil {
		return runtime.EmptyString, err
	}

	err = runtime.ValidateType(args[0], runtime.TypeString)

	if err != nil {
		return runtime.None, err
	}

	pathText := args[0].String()

	return runtime.NewString(path.Dir(pathText)), nil
}
