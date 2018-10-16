package html

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

func Fmt(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, core.MaxArgs)
	if err != nil {
		return values.None, err
	}

	err = core.ValidateType(args[0], core.StringType)
	if err != nil {
		return values.None, err
	}

	return values.None, nil
}
