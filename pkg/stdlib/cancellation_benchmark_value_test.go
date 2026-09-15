package stdlib_test

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// Expose traversal without measurable or membership fast paths.
type cancellationBenchmarkIterable struct {
	runtime.Value
	values *runtime.Array
}

func (value *cancellationBenchmarkIterable) Iterate(ctx context.Context) (runtime.Iterator, error) {
	return value.values.Iterate(ctx)
}
