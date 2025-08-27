package path

import (
	"context"
	"path"

	"github.com/MontFerret/ferret/pkg/runtime"
)

// EXT returns the extension of the last component of path.
// @param {String} path - The path.
// @return {String} - The extension of the last component of path.
func Ext(_ context.Context, args ...runtime.Value) (runtime.Value, error) {
	err := runtime.ValidateArgs(args, 1, 1)

	if err != nil {
		return runtime.EmptyString, err
	}

	err = runtime.ValidateType(args[0], runtime.TypeString)

	if err != nil {
		return runtime.EmptyString, err
	}

	pathText := args[0].String()

	return runtime.NewString(path.Ext(pathText)), nil
}
