package runtime

import "context"

type (
	// DispatchEvent represents an event that can be dispatched by a Dispatchable entity.
	DispatchEvent struct {
		Payload Value  `json:"payload"`
		Options Value  `json:"options"`
		Name    String `json:"name"`
	}

	// Dispatchable represents an entity that can dispatch events.
	// Dispatch is effectful and does not produce a runtime value.
	Dispatchable interface {
		Dispatch(ctx context.Context, event DispatchEvent) error
	}
)
