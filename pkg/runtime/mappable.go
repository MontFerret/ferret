package runtime

import "context"

type (
	// Mapper represents a function that maps a value and its key to a new value of type T.
	Mapper[T any] func(ctx context.Context, value, key Value) (T, error)
)
