package common

import (
	"context"
	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/runtime/core"
)

type Iterator struct {
	node drivers.HTMLElement
	pos  core.Int
}

func NewIterator(
	node drivers.HTMLElement,
) (core.Iterator, error) {
	if node == nil {
		return nil, core.Error(core.ErrMissedArgument, "result")
	}

	return &Iterator{node, 0}, nil
}

func (iterator *Iterator) HasNext(_ context.Context) (bool, error) {
	return iterator.node.Length() > int(iterator.pos), nil
}

func (iterator *Iterator) Next(ctx context.Context) (value core.Value, key core.Value, err error) {
	idx := iterator.pos
	val, err := iterator.node.GetChildNode(ctx, idx)

	if err != nil {
		return core.None, core.None, err
	}

	iterator.pos++

	return val, idx, nil
}
