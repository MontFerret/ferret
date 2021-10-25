package events

import (
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

type (
	// Event represents an event or value that an Observable can emit.
	Event interface {
		Value() core.Value
		Err() error
	}

	notify struct {
		value core.Value
		err   error
	}
)

func (n *notify) Value() core.Value {
	return n.value
}

func (n *notify) Err() error {
	return n.err
}

func WithValue(val core.Value) Event {
	return &notify{value: val}
}

func WithErr(err error) Event {
	return &notify{err: err, value: values.None}
}
