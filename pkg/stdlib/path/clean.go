package path

import (
	"context"
	"path"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

// Clean returns the shortest path name equivalent to path.
// @params path (String) - The path.
// @returns (String) - The shortest path name equivalent to path
func Clean(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, 1)

	if err != nil {
		return values.EmptyString, err
	}

	err = core.ValidateType(args[0], types.String)

	if err != nil {
		return values.None, err
	}

	pathText := args[0].String()

	return values.NewString(path.Clean(pathText)), nil
}
