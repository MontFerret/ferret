package events

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

type (
	strm chan Message

	msg struct {
		value core.Value
		err   error
	}
)

func New(source chan Message) Stream {
	return strm(source)
}

func (s strm) Close(_ context.Context) error {
	close(s)

	return nil
}

func (s strm) Read(ctx context.Context) <-chan Message {
	proxy := make(chan Message)

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

func (n *msg) Value() core.Value {
	return n.value
}

func (n *msg) Err() error {
	return n.err
}

func WithValue(val core.Value) Message {
	return &msg{value: val}
}

func WithErr(err error) Message {
	return &msg{err: err, value: values.None}
}
