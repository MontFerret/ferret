package runtime_test

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func compareValues(left, right runtime.Value) runtime.Ordering {
	result, err := runtime.CompareValues(context.Background(), left, right)
	if err != nil {
		panic(err)
	}

	return result
}
