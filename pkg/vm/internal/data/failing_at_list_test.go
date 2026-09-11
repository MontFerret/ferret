package data_test

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type failingAtList struct {
	runtime.List
	err error
}

func (l *failingAtList) At(context.Context, runtime.Int) (runtime.Value, error) {
	return runtime.None, l.err
}

func (l *failingAtList) New(ctx context.Context) (runtime.List, error) {
	base, err := l.List.New(ctx)
	if err != nil {
		return nil, err
	}

	result := *l
	result.List = base

	return &result, nil
}
