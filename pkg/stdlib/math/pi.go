package math

import (
	"context"
	"math"

	"github.com/MontFerret/ferret/pkg/runtime/core"
)

// PI returns Pi value.
// @return {Float} - Pi value.
func Pi(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 0, 0)

	if err != nil {
		return core.None, err
	}

	return core.NewFloat(math.Pi), nil
}
