package math

import (
	"context"
	"math"

	"github.com/MontFerret/ferret/pkg/runtime"
)

// PI returns Pi value.
// @return {Float} - Pi value.
func Pi(_ context.Context) (runtime.Value, error) {
	return runtime.NewFloat(math.Pi), nil
}
