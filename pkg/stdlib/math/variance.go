package math

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime"
)

func variance(ctx context.Context, input runtime.List, sample runtime.Int) (runtime.Float, error) {
	size, err := input.Length(ctx)

	if err != nil {
		return runtime.NaN(), err
	}

	if size == 0 {
		return runtime.NaN(), nil
	}

	m, err := mean(ctx, input)

	if err != nil {
		return runtime.NaN(), err
	}

	var variance runtime.Float

	err = input.ForEach(ctx, func(c context.Context, value runtime.Value, idx runtime.Int) (runtime.Boolean, error) {
		if err = runtime.AssertNumber(value); err != nil {
			return false, err
		}

		n := runtime.Float(toFloat(value))

		variance += (n - m) * (n - m)

		return true, nil
	})

	if err != nil {
		return runtime.NaN(), err
	}

	// When getting the mean of the squared differences
	// "sample" will allow us to know if it's a sample
	// or population and whether to subtract by one or not
	l := runtime.Float(size - (1 * sample))

	return variance / l, nil
}
