package events

import (
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

	// GenericSource represents a helper struct for generating custom event sources
	GenericSource struct {
		eventID ID
		stream  rpcc.Stream
		recv    func(stream rpcc.Stream) (interface{}, error)
	}
)

var (
	Error = New("error")
)

// NewSource create a new custom event source
// eventID - is a unique event ID
// stream - is a custom event stream
// recv - is a value conversion function
func NewSource(
	eventID ID,
	stream rpcc.Stream,
	recv func(stream rpcc.Stream) (interface{}, error),
) Source {
	return &GenericSource{eventID, stream, recv}
}

func (src *GenericSource) EventID() ID {
	return src.eventID
}

func (src *GenericSource) Ready() <-chan struct{} {
	return src.stream.Ready()
}

func (src *GenericSource) RecvMsg(m interface{}) error {
	return src.stream.RecvMsg(m)
}

func (src *GenericSource) Close() error {
	return src.stream.Close()
}

func (src *GenericSource) Recv() (Event, error) {
	data, err := src.recv(src.stream)

	if err != nil {
		return Event{}, err
	}

	return Event{
		ID:   src.eventID,
		Data: data,
	}, nil
}
