package math

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime"
)

func mean(ctx context.Context, input runtime.List) (runtime.Float, error) {
	size, err := input.Length(ctx)

	if err != nil {
		return 0, err
	}

	if size == 0 {
		return runtime.NaN(), nil
	}

	var sum float64

	err = input.ForEach(ctx, func(c context.Context, value runtime.Value, idx runtime.Int) (runtime.Boolean, error) {
		err = runtime.AssertNumber(value)

		if err != nil {
			return false, nil
		}

		sum += toFloat(value)

		return true, nil
	})

	if err != nil {
		return 0, err
	}

	return runtime.NewFloat(sum / float64(size)), nil
}
