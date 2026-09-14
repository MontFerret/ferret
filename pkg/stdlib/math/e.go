package math

import (
	"context"
	"math"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// e returns Euler's number, the base of the natural logarithm.
// @return {Float} Euler's number, approximately 2.718281828459045.
func E(_ context.Context) (runtime.Value, error) {
	return runtime.NewFloat(math.E), nil
}
