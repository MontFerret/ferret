package core

import "context"

type (
	Indexed interface {
		Get(ctx context.Context, idx int) (Value, error)
	}

	Keyed interface {
		Get(ctx context.Context, key string) (Value, error)
	}

	// Collection represents a collection of values.
	// Generic interface for all collection-like structures.
	Collection interface {
		Value
		Comparable
		Measurable
		Iterable

		IsEmpty(ctx context.Context) (bool, error)
		Clear(ctx context.Context) error
	}

	// List represents a list of values.
	// Generic interface for all list-like structures.
	List interface {
		Collection
		Indexed

		Contains(ctx context.Context, value Value) (bool, error)
		Add(ctx context.Context, value Value) error
		Remove(ctx context.Context, value Value) error
		RemoveAt(ctx context.Context, idx int) (Value, error)
		Swap(ctx context.Context, i, j int) error
	}

	// Set represents a set of values.
	// Generic interface for all set-like structures.
	Set interface {
		Collection

		Contains(ctx context.Context, value Value) (bool, error)
		Add(ctx context.Context, value Value) error
		Remove(ctx context.Context, value Value) error
	}

	// Map represents a dictionary of values.
	// Generic interface for all dictionary-like structures.
	Map interface {
		Collection
		Keyed

		ContainsKey(ctx context.Context, key string) (bool, error)
		ContainsValue(ctx context.Context, value Value) (bool, error)
		Put(ctx context.Context, key string, value Value) error
		Remove(ctx context.Context, key string) error
	}
)
