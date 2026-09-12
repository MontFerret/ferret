package math

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func mean(ctx context.Context, input runtime.List) (runtime.Float, runtime.Int, error) {
	sum, count, err := sumNumbers(ctx, input)
	if err != nil {
		return runtime.NaN(), 0, err
	}

	if count == 0 {
		return runtime.NaN(), 0, nil
	}

	return runtime.Float(sum / float64(count)), count, nil
}
