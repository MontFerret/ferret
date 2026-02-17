package runtime

import "context"

type (
	// DispatchEvent represents an event that can be dispatched by a Dispatchable entity.
	DispatchEvent struct {
		Name    String
		Payload Value
		Options Value
	}

	// Dispatchable represents an entity that can dispatch events.
	Dispatchable interface {
		Dispatch(ctx context.Context, event DispatchEvent) (Value, error)
	}
)
