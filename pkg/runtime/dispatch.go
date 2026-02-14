package runtime

import "context"

type (
	DispatchEvent struct {
		Name    String
		Payload Value
		Options Value
	}

	Dispatcher interface {
		Dispatch(ctx context.Context, event DispatchEvent) (Value, error)
	}
)
