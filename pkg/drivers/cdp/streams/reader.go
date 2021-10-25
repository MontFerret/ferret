package streams

import (
	"context"
	"github.com/mafredri/cdp/rpcc"

	"github.com/MontFerret/ferret/pkg/runtime/events"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

type Reader struct {
	decoder Decoder
}

func NewReader(decoder Decoder) *Reader {
	return &Reader{decoder}
}

func (reader *Reader) Read(ctx context.Context, stream rpcc.Stream) <-chan events.Event {
	ch := make(chan events.Event)

	go func() {
		defer func() {
			_ = stream.Close()

			close(ch)
		}()

		for {
			select {
			case <-ctx.Done():
				return
			case <-stream.Ready():
				val, err := reader.decoder(stream)

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
