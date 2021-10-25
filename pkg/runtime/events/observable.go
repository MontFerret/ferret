package events

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

type (
	// Subscription represents an event subscription object that contains target event name
	// and optional event options.
	Subscription struct {
		EventName string
		Options   *values.Object
	}

	// Observable represents an interface of
	// complex types that can emit events.
	Observable interface {
		Subscribe(ctx context.Context, subscription Subscription) (<-chan Event, error)
	}
)
