package events

import (
	"context"
	"github.com/mafredri/cdp/rpcc"
)

type (
	// ID represents a unique event ID
	ID int

	// Event represents a system event that is returned from an event source
	Event struct {
		ID   ID
		Data interface{}
	}

	// Source represents a custom source of system events
	Source interface {
		rpcc.Stream
		Recv() (Event, error)
	}

	// SourceFactory represents a function that creates a new instance of Source.
	SourceFactory func(ctx context.Context) (Source, error)

	StreamFactory func(ctx context.Context) (rpcc.Stream, error)

	StreamDecoder func(stream rpcc.Stream) (interface{}, error)

	// StreamSource represents a helper struct for generating custom event sources
	StreamSource struct {
		eventID ID
		stream  rpcc.Stream
		decoder StreamDecoder
	}
)

var (
	Error = New("error")
)

// NewStreamSource create a new custom event source based on rpcc.Stream
// eventID - is a unique event ID
// stream - is a custom event stream
// decoder - is a value conversion function
func NewStreamSource(
	eventID ID,
	stream rpcc.Stream,
	decoder StreamDecoder,
) Source {
	return &StreamSource{eventID, stream, decoder}
}

func (src *StreamSource) ID() ID {
	return src.eventID
}

func (src *StreamSource) Ready() <-chan struct{} {
	return src.stream.Ready()
}

func (src *StreamSource) RecvMsg(m interface{}) error {
	return src.stream.RecvMsg(m)
}

func (src *StreamSource) Close() error {
	return src.stream.Close()
}

func (src *StreamSource) Recv() (Event, error) {
	data, err := src.decoder(src.stream)

	if err != nil {
		return Event{}, err
	}

	return Event{
		ID:   src.eventID,
		Data: data,
	}, nil
}

func NewStreamSourceFactory(eventID ID, factory StreamFactory, receiver StreamDecoder) SourceFactory {
	return func(ctx context.Context) (Source, error) {
		stream, err := factory(ctx)

		if err != nil {
			return nil, err
		}

		return NewStreamSource(eventID, stream, receiver), nil
	}
}
