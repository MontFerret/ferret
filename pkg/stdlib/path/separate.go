package path

import (
	"context"
	"path"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

// SEPARATE separates the path into a directory and filename component.
// @param {String} path - The path
// @return {Any[]} - First item is a directory component, and second is a filename component.
func Separate(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, 1)

	if err != nil {
		return values.None, err
	}

	err = core.ValidateType(args[0], types.String)

	if err != nil {
		return values.None, err
	}

	pattern, name := path.Split(args[0].String())

	arr := values.NewArrayWith(values.NewString(pattern), values.NewString(name))

	return arr, nil
}
