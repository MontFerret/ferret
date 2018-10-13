package math

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"math/rand"
)

/*
 * Return a pseudo-random number between 0 and 1.
 * @returns (Float) - A number greater than 0 and less than 1.
 */
func Rand(_ context.Context, _ ...core.Value) (core.Value, error) {
	return values.NewFloat(rand.Float64()), nil
}
