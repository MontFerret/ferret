package random

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// bool draws a non-cryptographic pseudo-random Boolean from the Session source.
// @return {Boolean} True or false with equal probability, consuming the shared random sequence.
func Bool(ctx context.Context) (runtime.Value, error) {
	src, err := source(ctx)
	if err != nil {
		return runtime.None, err
	}

	return runtime.Boolean(src.Bool()), nil
}
