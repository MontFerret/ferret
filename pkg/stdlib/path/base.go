package path

import (
	"context"
	"path"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

// Base returns the last component of the path.
// or the path itself if it does not contain any directory separators.
// @params path (String) - The path.
// @returns (String) - The last component of the path.
func Base(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, 1)

	if err != nil {
		return values.EmptyString, err
	}

	err = core.ValidateType(args[0], types.String)

	if err != nil {
		return values.None, err
	}

	pathText := args[0].String()

	return values.NewString(path.Base(pathText)), nil
}
