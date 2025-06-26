package types

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime"

	"github.com/MontFerret/ferret/pkg/runtime/core"
)

// TO_DATETIME takes an input value of any type and converts it into the appropriate date time value.
// @param {Any} value - Input value of arbitrary type.
// @return {DateTime} - Parsed date time.
func ToDateTime(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, 1)

	if err != nil {
		return core.None, err
	}

	return runtime.ParseDateTime(args[0].String())
}
