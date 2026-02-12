package runtime

import (
	"io"
)

type (
	ValueStorage interface {
		Set(ctx Context, value Value) error
		Get(ctx Context) (Value, error)
	}

	MapStorage interface {
		Iterable
		Measurable
		KeyReadable

		Set(ctx Context, key, value Value) error
	}

	ListStorage interface {
		Iterable
		Measurable
		IndexReadable
		Sortable

		Add(ctx Context, value Value) error
		AddKV(ctx Context, key, value Value) error
	}

	StorageManager interface {
		io.Closer
		ValueStorage(ctx Context, name String) (ValueStorage, error)
		ListStorage(ctx Context, name String) (ListStorage, error)
		MapStorage(ctx Context, name String) (MapStorage, error)
	}

	StateManagerFactory interface {
		New(ctx Context) (StorageManager, error)
	}
)
