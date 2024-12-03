package events

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/values"

	"github.com/MontFerret/ferret/pkg/runtime/core"
)

type Iterator struct {
	messages <-chan Message
	message  Message
}

func NewIterator(ch <-chan Message) core.Iterator {
	return &Iterator{ch, nil}
}

func (iter *Iterator) HasNext(ctx context.Context) (bool, error) {
	select {
	case evt, ok := <-iter.messages:
		if !ok {
			return false, nil
		}

		iter.message = evt

		return true, nil
	case <-ctx.Done():
		return false, ctx.Err()
	}
}

func (iter *Iterator) Next(ctx context.Context) (value core.Value, key core.Value, err error) {
	if iter.message != nil {
		if err := iter.message.Err(); err != nil {
			return values.None, values.None, err
		}

		return iter.message.Value(), values.None, nil
	}

	return values.None, values.None, core.ErrNoMoreData
}
