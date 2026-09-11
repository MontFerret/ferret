package runtime_test

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type copyContractList struct {
	runtime.List
	copied runtime.Value
	calls  int
}

func (list *copyContractList) Copy() runtime.Value {
	list.calls++

	return list.copied
}

func (list *copyContractList) New(ctx context.Context) (runtime.List, error) {
	base, err := list.List.New(ctx)
	if err != nil {
		return nil, err
	}

	result := *list
	result.List = base

	return &result, nil
}
