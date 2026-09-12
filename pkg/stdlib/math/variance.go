package math

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func variance(ctx context.Context, input runtime.List, sample bool) (runtime.Float, error) {
	m, count, err := mean(ctx, input)
	if err != nil {
		return runtime.NaN(), err
	}

	if count == 0 || (sample && count < 2) {
		return runtime.NaN(), nil
	}

	var squaredDifferences float64

	_, err = forEachNumber(ctx, input, func(value runtime.Value, _ runtime.Int) {
		difference := toFloat(value) - float64(m)
		squaredDifferences += difference * difference
	})
	if err != nil {
		return runtime.NaN(), err
	}

	if sample {
		count--
	}

	return runtime.Float(squaredDifferences / float64(count)), nil
}
