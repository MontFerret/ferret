package events

import (
	"context"
	"github.com/mafredri/cdp/rpcc"
)

type (
	ID int

	Event struct {
		ID   ID
		Data interface{}
	}

	Handler func(ctx context.Context, message interface{})

	Listener struct {
		EventID ID
		Handler Handler
	}

	Source interface {
		rpcc.Stream
		Recv() (Event, error)
	}

	GenericSource struct {
		eventID ID
		stream  rpcc.Stream
		recv    func() (interface{}, error)
	}
)

const (
	//revive:disable-next-line:var-declaration
	IDAny = ID(iota)
	IDError
	IDLoad
	IDReload
	IDAttrModified
	IDAttrRemoved
	IDChildNodeCountUpdated
	IDChildNodeInserted
	IDChildNodeRemoved
)

func NewSource(
	eventID ID,
	stream rpcc.Stream,
	recv func() (interface{}, error),
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
	data, err := src.recv()

	if err != nil {
		return Event{}, err
	}

	return Event{
		ID:   src.eventID,
		Data: data,
	}, nil
}
