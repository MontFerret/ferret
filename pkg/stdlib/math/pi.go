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
func Pi(_ context.Context, _ ...core.Value) (core.Value, error) {
	return values.NewFloat(math.Pi), nil
}
