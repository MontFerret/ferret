package runtime

import (
	"context"
	"io"
)

type (
	ValueStorage interface {
		Set(ctx context.Context, value Value) error
		Get(ctx context.Context) (Value, error)
	}

	MapStorage interface {
		Iterable
		Measurable
		Keyed

		Set(ctx context.Context, key, value Value) error
	}

	ListStorage interface {
		Iterable
		Measurable
		Indexed

		Add(ctx context.Context, value Value) error
	}

	StorageManager interface {
		io.Closer
		ValueStorage(ctx context.Context, name String) (ValueStorage, error)
		ListStorage(ctx context.Context, name String) (ListStorage, error)
		MapStorage(ctx context.Context, name String) (MapStorage, error)
	}

	StateManagerFactory interface {
		New(ctx context.Context) (StorageManager, error)
	}
)
