package math

import (
	"context"
	"math"

	"github.com/MontFerret/ferret/pkg/runtime"
)

// PI returns Pi value.
// @return {Float} - Pi value.
func Pi(_ context.Context, args ...runtime.Value) (runtime.Value, error) {
	err := runtime.ValidateArgs(args, 0, 0)

	if err != nil {
		return runtime.None, err
	}

	return runtime.NewFloat(math.Pi), nil
}
