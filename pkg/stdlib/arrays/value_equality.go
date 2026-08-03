package arrays

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func findEqualValue(ctx context.Context, candidates []runtime.Value, target runtime.Value) (int, error) {
	for idx, candidate := range candidates {
		equal, err := runtime.EqualValues(ctx, candidate, target)
		if err != nil {
			return -1, err
		}

		if equal {
			return idx, nil
		}
	}

	return -1, nil
}
