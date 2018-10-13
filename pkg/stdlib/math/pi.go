package math

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"math"
)

func Pi(_ context.Context, _ ...core.Value) (core.Value, error) {
	return values.NewFloat(math.Pi), nil
}
