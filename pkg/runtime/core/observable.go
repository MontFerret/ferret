package core

import (
	"context"
)

// Observable represents an interface of
// complex types that can have event subscribers.
type Observable interface {
	Subscribe(ctx context.Context, eventName string) (<-chan struct{}, error)
}
