package events

import (
	"github.com/mafredri/cdp/rpcc"
)

type (
	Type int

	Event struct {
		Type Type
		Data interface{}
	}

	Listener struct {
		Event   Type
		Handler EventHandler
	}

	Source interface {
		rpcc.Stream
		Recv() (Event, error)
	}

	GenericSource struct {
		evtType Type
		stream  rpcc.Stream
		recv    func() (interface{}, error)
	}
)

const (
	//revive:disable-next-line:var-declaration
	EventTypeAny = Type(iota)
	EventTypeError
	EventTypeLoad
	EventTypeReload
	EventTypeAttrModified
	EventTypeAttrRemoved
	EventTypeChildNodeCountUpdated
	EventTypeChildNodeInserted
	EventTypeChildNodeRemoved
)

func NewSource(
	evtType Type,
	stream rpcc.Stream,
	recv func() (interface{}, error),
) GenericSource {
	return GenericSource{evtType, stream, recv}
}

func (src GenericSource) Ready() <-chan struct{} {
	return src.stream.Ready()
}

func (src GenericSource) RecvMsg(m interface{}) error {
	return src.stream.RecvMsg(m)
}

func (src GenericSource) Close() error {
	return src.stream.Close()
}

func (src GenericSource) Recv() (Event, error) {
	data, err := src.recv()

	if err != nil {
		return Event{}, err
	}

	return Event{
		Type: src.evtType,
		Data: data,
	}, nil
}
