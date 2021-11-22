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

	// Message represents an event message that an Observable can emit.
	Message interface {
		Value() core.Value
		Err() error
	}

	// Stream represents an event stream that produces target event objects.
	Stream interface {
		Close(ctx context.Context) error
		Read(ctx context.Context) <-chan Message
	}

	// Observable represents an interface of
	// complex types that returns stream of events.
	Observable interface {
		Subscribe(ctx context.Context, subscription Subscription) (Stream, error)
	}
)
