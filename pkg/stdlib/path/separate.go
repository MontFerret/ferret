package path

import (
	"path"

	"github.com/MontFerret/ferret/pkg/runtime"
)

// SEPARATE separates the path into a directory and filename component.
// @param {String} path - The path
// @return {Any[]} - First item is a directory component, and second is a filename component.
func Separate(ctx runtime.Context, args ...runtime.Value) (runtime.Value, error) {
	err := runtime.ValidateArgs(args, 1, 1)

	if err != nil {
		return runtime.None, err
	}

	err = runtime.ValidateType(args[0], runtime.TypeString)

	if err != nil {
		return runtime.None, err
	}

	pattern, name := path.Split(args[0].String())

	arr := ctx.Alloc().Array(2)
	_ = arr.Append(ctx, runtime.NewString(pattern))
	_ = arr.Append(ctx, runtime.NewString(name))

	return arr, nil
}
