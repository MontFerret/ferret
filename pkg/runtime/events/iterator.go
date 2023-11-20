package events

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

type Iterator struct {
	messages <-chan Message
	message  Message
}

func NewIterator(ch <-chan Message) core.Iterator {
	return &Iterator{ch, nil}
}

func (e *Iterator) HasNext(ctx context.Context) (bool, error) {
	select {
	case evt, ok := <-e.messages:
		if !ok {
			return false, nil
		}

		e.message = evt

		return true, nil
	case <-ctx.Done():
		return false, ctx.Err()
	}
}

func (e *Iterator) Next(_ context.Context) (value core.Value, key core.Value, err error) {
	if e.message != nil {
		if err := e.message.Err(); err != nil {
			return values.None, values.None, err
		}

		return e.message.Value(), values.None, nil
	}

	return values.None, values.None, core.ErrNoMoreData
}
