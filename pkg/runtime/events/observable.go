package events

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

type (
	// Event represents an event object that contains either an optional event data
	// or error that occurred during an event
	Event struct {
		Data core.Value
		Err  error
	}

	// Subscription represents an event subscription object that contains target event name
	// and optional event options.
	Subscription struct {
		EventName string
		Options   *values.Object
	}

	// Observable represents an interface of
	// complex types that can have event subscribers.
	Observable interface {
		Subscribe(ctx context.Context, subscription Subscription) (<-chan Event, error)
	}
)
