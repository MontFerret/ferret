package math

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"math"
)

/*
 * Returns Pi value.
 * @returns (Float) - Pi value.
 */
func Pi(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 0, 0)

	if err != nil {
		return values.None, err
	}

	return values.NewFloat(math.Pi), nil
}
