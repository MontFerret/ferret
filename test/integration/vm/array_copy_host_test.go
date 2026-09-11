package vm_test

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type invalidArrayCopy struct {
	runtime.List
}

func (list *invalidArrayCopy) Copy() runtime.Value {
	return runtime.True
}

func (list *invalidArrayCopy) New(ctx context.Context) (runtime.List, error) {
	base, err := list.List.New(ctx)
	if err != nil {
		return nil, err
	}

	result := *list
	result.List = base

	return &result, nil
}
