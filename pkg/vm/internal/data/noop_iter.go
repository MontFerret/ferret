package data

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime"
)

type noopIter struct{}

func (n noopIter) HasNext(ctx runtime.Context) (bool, error) {
	return false, nil
}

func (n noopIter) Next(ctx runtime.Context) (value runtime.Value, key runtime.Value, err error) {
	return runtime.None, runtime.None, nil
}
