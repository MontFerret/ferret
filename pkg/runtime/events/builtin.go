package events

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

type (
	strm chan Event

	evt struct {
		value core.Value
		err   error
	}
)

func New(source chan Event) Stream {
	return strm(source)
}

func (s strm) Close(_ context.Context) error {
	close(s)

	return nil
}

func (s strm) Read(ctx context.Context) <-chan Event {
	proxy := make(chan Event)

	go func() {
		defer close(proxy)

		for {
			select {
			case <-ctx.Done():
				return
			case evt := <-s:
				if ctx.Err() != nil {
					return
				}

				proxy <- evt
			}
		}
	}()

	return s
}

func (n *evt) Value() core.Value {
	return n.value
}

func (n *evt) Err() error {
	return n.err
}

func WithValue(val core.Value) Event {
	return &evt{value: val}
}

func WithErr(err error) Event {
	return &evt{err: err, value: values.None}
}
