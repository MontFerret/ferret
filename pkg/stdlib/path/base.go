package path

import (
	"context"
	"path"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

// BASE returns the last component of the path or the path itself if it does not contain any directory separators.
// @param {String} path - The path.
// @return {String} - The last component of the path.
func Base(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, 1)

	if err != nil {
		return core.EmptyString, err
	}

	err = core.ValidateType(args[0], types.String)

	if err != nil {
		return core.None, err
	}

	pathText := args[0].String()

	return core.NewString(path.Base(pathText)), nil
}
