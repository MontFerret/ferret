package math

import (
	"context"
	"math/rand"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

// Rand return a pseudo-random number between 0 and 1.
// @returns (Float) - A number greater than 0 and less than 1.
func Rand(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 0, 0)

	if err != nil {
		return values.None, err
	}

	return values.NewFloat(rand.Float64()), nil
}
