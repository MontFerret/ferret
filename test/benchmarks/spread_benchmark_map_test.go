package benchmarks_test

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type benchmarkMap struct {
	runtime.Map
}

func (source benchmarkMap) New(ctx context.Context) (runtime.Map, error) {
	base, err := source.Map.New(ctx)
	if err != nil {
		return nil, err
	}

	result := source
	result.Map = base

	return result, nil
}
