package math

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func extremum(ctx context.Context, source runtime.List, minimum bool, policy numberPolicy) (runtime.Value, error) {
	var result float64

	counts, err := forEachNumber(ctx, source, policy, func(value runtime.Value, index runtime.Int) {
		number := toFloat(value)
		if index == 0 || (minimum && number < result) || (!minimum && number > result) {
			result = number
		}
	})
	if err != nil {
		return runtime.None, err
	}

	if counts.numeric == 0 {
		return runtime.None, nil
	}

	return runtime.Float(result), nil
}
