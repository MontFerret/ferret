package events

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

type (
	// Subscription represents an event subscription object that contains target event name
	// and optional event options.
	Subscription struct {
		EventName string
		Options   *values.Object
	}

	// Event represents an event or value that an Observable can emit.
	Event interface {
		Value() core.Value
		Err() error
	}

	// Stream represents an event stream that produces target event objects.
	Stream interface {
		Close(ctx context.Context) error
		Read(ctx context.Context) <-chan Event
	}

	// Observable represents an interface of
	// complex types that returns stream of events.
	Observable interface {
		Subscribe(ctx context.Context, subscription Subscription) (Stream, error)
	}
)
