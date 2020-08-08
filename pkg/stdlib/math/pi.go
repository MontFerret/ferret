package math

import (
	"context"
	"math"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

// PI returns Pi value.
// @return {Float} - Pi value.
func Pi(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 0, 0)

	if err != nil {
		return values.None, err
	}

	return values.NewFloat(math.Pi), nil
}
