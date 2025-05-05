package internal

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime"
)

type noopIter struct{}

func (n noopIter) HasNext(_ context.Context) (bool, error) {
	return false, nil
}

func (n noopIter) Next(_ context.Context) (value runtime.Value, key runtime.Value, err error) {
	return runtime.None, runtime.None, nil
}
