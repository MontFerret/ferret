package path

import (
	"context"
	"path"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// SEPARATE separates the path into a directory and filename component.
// @param {String} path - The path
// @return {Any[]} - First item is a directory component, and second is a filename component.
func Separate(_ context.Context, arg runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateType(arg, runtime.TypeString); err != nil {
		return runtime.None, err
	}

	pattern, name := path.Split(arg.String())

	arr := runtime.NewArrayWith(runtime.NewString(pattern), runtime.NewString(name))

	return arr, nil
}
