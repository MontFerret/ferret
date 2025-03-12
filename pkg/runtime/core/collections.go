package core

import "context"

type (
	Indexed interface {
		Get(ctx context.Context, idx Int) (Value, error)
	}

	IndexedPredicate = func(ctx context.Context, value Value, idx Int) (bool, error)

	Keyed interface {
		Get(ctx context.Context, key Value) (Value, error)
	}

	KeyedPredicate = func(ctx context.Context, value, key Value) (bool, error)

	// Collection represents a collection of values.
	// Generic interface for all collection-like structures.
	Collection interface {
		Value
		Comparable
		Measurable
		Cloneable
		Iterable

		IsEmpty(ctx context.Context) (bool, error)
		Clear(ctx context.Context) error
	}

	// List represents a items of values.
	// Generic interface for all items-like structures.
	List interface {
		Collection
		Indexed

		ForEach(ctx context.Context, predicate IndexedPredicate) error

		Find(ctx context.Context, predicate IndexedPredicate) (List, error)
		FindOne(ctx context.Context, predicate IndexedPredicate) (Value, Boolean, error)
		Contains(ctx context.Context, value Value) (Boolean, error)
		IndexOf(ctx context.Context, value Value) (Int, error)
		First(context.Context) (Value, error)
		Last(context.Context) (Value, error)
		Slice(ctx context.Context, start, end Int) (List, error)
		Sort(ctx context.Context, ascending Boolean) (List, error)
		SortWith(ctx context.Context, comparator Comparator) (List, error)

		Add(ctx context.Context, value Value) error
		Set(ctx context.Context, idx Int, value Value) error
		Insert(ctx context.Context, idx Int, value Value) error
		Remove(ctx context.Context, value Value) error
		RemoveAt(ctx context.Context, idx Int) (Value, error)
		Swap(ctx context.Context, i, j Int) error
	}

	// Set represents a set of values.
	// Generic interface for all set-like structures.
	Set interface {
		Collection

		Contains(ctx context.Context, value Value) (Boolean, error)
		Add(ctx context.Context, value Value) error
		Remove(ctx context.Context, value Value) error
	}

	// Map represents a dictionary of values.
	// Generic interface for all dictionary-like structures.
	Map interface {
		Collection
		Keyed

		ForEach(ctx context.Context, predicate KeyedPredicate) error

		Keys(context.Context) ([]Value, error)
		Values(context.Context) ([]Value, error)

		Find(ctx context.Context, predicate KeyedPredicate) (List, error)
		FindOne(ctx context.Context, predicate KeyedPredicate) (Value, Boolean, error)
		ContainsKey(ctx context.Context, key Value) (Boolean, error)
		ContainsValue(ctx context.Context, value Value) (Boolean, error)

		Set(ctx context.Context, key Value, value Value) error
		Remove(ctx context.Context, key Value) error
	}
)
