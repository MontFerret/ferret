package internal

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
)

type noopIter struct{}

func (n noopIter) HasNext(_ context.Context) (bool, error) {
	return false, nil
}

func (n noopIter) Next(_ context.Context) (value core.Value, key core.Value, err error) {
	return core.None, core.None, nil
}
