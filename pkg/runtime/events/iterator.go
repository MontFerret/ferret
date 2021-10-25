package events

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

type Iterator struct {
	ch <-chan Event
}

func NewIterator(ch <-chan Event) core.Iterator {
	return &Iterator{ch}
}

func (e *Iterator) Next(ctx context.Context) (value core.Value, key core.Value, err error) {
	select {
	case evt, ok := <-e.ch:
		if !ok {
			return values.None, values.None, core.ErrNoMoreData
		}

		if err := evt.Err(); err != nil {
			return values.None, values.None, err
		}

		return evt.Value(), values.None, nil
	case <-ctx.Done():
		return values.None, values.None, ctx.Err()
	}
}
