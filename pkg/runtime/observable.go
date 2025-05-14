package runtime

import (
	"context"
	"io"
)

const DefaultStreamTimeout = 5000

type (
	// Subscription represents an event subscription object that contains target event name
	// and optional event options.
	Subscription struct {
		EventName String
		Options   Map
	}

	// Message represents an event message that an Observable can emit.
	Message interface {
		Value() Value
		Err() error
	}

	// Stream represents an event stream that produces target event objects.
	Stream interface {
		io.Closer
		Read(ctx context.Context) <-chan Message
	}

	// Observable represents an interface of
	// complex types that returns stream of events.
	Observable interface {
		Subscribe(ctx context.Context, subscription Subscription) (Stream, error)
	}
)
