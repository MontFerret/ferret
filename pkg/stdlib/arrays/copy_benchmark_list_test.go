package arrays_test

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type copyBenchmarkList struct {
	runtime.List
}

func (source copyBenchmarkList) New(ctx context.Context) (runtime.List, error) {
	base, err := source.List.New(ctx)
	if err != nil {
		return nil, err
	}

	result := source
	result.List = base

	return result, nil
}
