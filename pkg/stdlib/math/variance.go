package math

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func variance(ctx context.Context, input runtime.List, sample bool, policy numberPolicy) (runtime.Float, error) {
	var mean, squaredDifferences float64

	// Welford's recurrence keeps the source single-pass and avoids subtracting
	// large raw sums of squares when the variance is small relative to the mean.
	counts, err := forEachNumber(ctx, input, policy, func(value runtime.Value, index runtime.Int) {
		number := toFloat(value)
		difference := number - mean
		mean += difference / float64(index+1)
		squaredDifferences += difference * (number - mean)
	})
	if err != nil {
		return runtime.NaN(), err
	}

	count := counts.numeric
	if count == 0 || (sample && count < 2) {
		return runtime.NaN(), nil
	}

	if sample {
		count--
	}

	return runtime.Float(squaredDifferences / float64(count)), nil
}
