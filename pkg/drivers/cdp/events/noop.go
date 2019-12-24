package events

import (
	"github.com/MontFerret/ferret/pkg/runtime/core"
)

type noopEvent struct {
	c chan struct{}
}

func newNoopSource() Source {
	return noopEvent{
		c: make(chan struct{}),
	}
}

func (n noopEvent) Ready() <-chan struct{} {
	return n.c
}

func (n noopEvent) RecvMsg(_ interface{}) error {
	return core.ErrNotSupported
}

func (n noopEvent) Close() error {
	return nil
}

func (n noopEvent) Recv() (Event, error) {
	return Event{}, core.ErrNotSupported
}
