package internal

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
)

type noopIter struct{}

var NoopIter = &noopIter{}

func (it *noopIter) HasNext(ctx context.Context) (bool, error) {
	return false, nil
}

func (it *noopIter) Next(_ context.Context) error {
	return nil
}

func (it *noopIter) Value() core.Value {
	return core.None
}

func (it *noopIter) Key() core.Value {
	return core.None
}
