package path

import (
	"context"
	"path"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

// Match reports whether name matches the pattern.
// @param pattern (String) - The pattern.
// @param name (String) - The name.
// @returns (Boolean) - True if the name mathes the pattern.
func Match(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 2, 2)

	if err != nil {
		return values.False, err
	}

	err = core.ValidateType(args[0], types.String)

	if err != nil {
		return values.False, err
	}

	err = core.ValidateType(args[1], types.String)

	if err != nil {
		return values.False, err
	}

	pattern := args[0].String()
	name := args[1].String()

	matched, err := path.Match(pattern, name)

	if err != nil {
		return values.False, core.Error(err, "match")
	}

	return values.NewBoolean(matched), nil
}
