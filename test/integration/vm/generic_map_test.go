package vm_test

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type genericMap struct {
	runtime.Map
}

func (genericMap) Type() runtime.Type {
	return runtime.TypeMap
}

func (source genericMap) New(ctx context.Context) (runtime.Map, error) {
	base, err := source.Map.New(ctx)
	if err != nil {
		return nil, err
	}

	result := source
	result.Map = base

	return result, nil
}
