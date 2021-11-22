package events

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/events"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/mafredri/cdp/rpcc"
)

type (
	Decoder func(ctx context.Context, stream rpcc.Stream) (core.Value, error)

	Factory func(ctx context.Context) (rpcc.Stream, error)

	EventStream struct {
		stream  rpcc.Stream
		decoder Decoder
	}
)

func NewEventStream(stream rpcc.Stream, decoder Decoder) events.Stream {
	return &EventStream{stream, decoder}
}

func (e *EventStream) Close(_ context.Context) error {
	return e.stream.Close()
}

func (e *EventStream) Read(ctx context.Context) <-chan events.Message {
	ch := make(chan events.Message)

	go func() {
		defer close(ch)

		for {
			select {
			case <-ctx.Done():
				return
			case <-e.stream.Ready():
				val, err := e.decoder(ctx, e.stream)

				if err != nil {
					ch <- events.WithErr(err)

					return
				}

				if val != nil && val != values.None {
					ch <- events.WithValue(val)
				}
			}
		}
	}()

	return ch
}
