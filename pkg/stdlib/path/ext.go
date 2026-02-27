package path

import (
	"context"
	"path"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// EXT returns the extension of the last component of path.
// @param {String} path - The path.
// @return {String} - The extension of the last component of path.
func Ext(_ context.Context, arg runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateType(arg, runtime.TypeString); err != nil {
		return runtime.EmptyString, err
	}

	pathText := arg.String()

	return runtime.NewString(path.Ext(pathText)), nil
}
