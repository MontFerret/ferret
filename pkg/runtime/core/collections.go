package core

import "context"

type (
	Indexed interface {
		GetByIndex(ctx context.Context, idx int) (Value, error)
	}

	Keyed interface {
		GetByKey(ctx context.Context, key string) (Value, error)
	}
)
