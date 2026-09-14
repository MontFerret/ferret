package math

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func mean(ctx context.Context, input runtime.List) (runtime.Float, runtime.Int, error) {
	sum, counts, err := sumNumbers(ctx, input)
	if err != nil {
		return runtime.NaN(), 0, err
	}

	if counts.numeric == 0 {
		return runtime.ZeroFloat, 0, nil
	}

	return runtime.Float(sum / float64(counts.numeric)), counts.numeric, nil
}
