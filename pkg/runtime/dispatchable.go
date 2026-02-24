package runtime

import "context"

type (
	// DispatchEvent represents an event that can be dispatched by a Dispatchable entity.
	DispatchEvent struct {
		Name    String `json:"name"`
		Payload Value  `json:"payload"`
		Options Value  `json:"options"`
	}

	// Dispatchable represents an entity that can dispatch events.
	Dispatchable interface {
		Dispatch(ctx context.Context, event DispatchEvent) (Value, error)
	}
)
