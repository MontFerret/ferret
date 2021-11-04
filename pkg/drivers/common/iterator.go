package common

import (
	"context"

	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

type Iterator struct {
	node drivers.HTMLElement
	pos  values.Int
}

func NewIterator(
	node drivers.HTMLElement,
) (core.Iterator, error) {
	if node == nil {
		return nil, core.Error(core.ErrMissedArgument, "result")
	}

	return &Iterator{node, 0}, nil
}

func (iterator *Iterator) Next(ctx context.Context) (value core.Value, key core.Value, err error) {
	if iterator.node.Length() > iterator.pos {
		idx := iterator.pos
		val, err := iterator.node.GetChildNode(ctx, idx)

		if err != nil {
			return values.None, values.None, err
		}

		iterator.pos++

		return val, idx, nil
	}

	return values.None, values.None, core.ErrNoMoreData
}
