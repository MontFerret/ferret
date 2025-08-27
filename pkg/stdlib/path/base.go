package path

import (
	"context"
	"path"

	"github.com/MontFerret/ferret/pkg/runtime"
)

// BASE returns the last component of the path or the path itself if it does not contain any directory separators.
// @param {String} path - The path.
// @return {String} - The last component of the path.
func Base(_ context.Context, args ...runtime.Value) (runtime.Value, error) {
	err := runtime.ValidateArgs(args, 1, 1)

	if err != nil {
		return runtime.EmptyString, err
	}

	err = runtime.ValidateType(args[0], runtime.TypeString)

	if err != nil {
		return runtime.None, err
	}

	pathText := args[0].String()

	return runtime.NewString(path.Base(pathText)), nil
}
