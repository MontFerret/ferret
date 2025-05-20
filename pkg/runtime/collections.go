package runtime

import "context"

type (
	Predicate = func(ctx context.Context, value, idx Value) (Boolean, error)

	Indexed interface {
		Get(ctx context.Context, idx Int) (Value, error)
	}

	IndexedPredicate = func(ctx context.Context, value Value, idx Int) (Boolean, error)

	Keyed interface {
		Get(ctx context.Context, key Value) (Value, error)
	}

	KeyedPredicate = func(ctx context.Context, value, key Value) (Boolean, error)

	// Collection represents a collection of values.
	// Generic interface for all collection-like structures.
	Collection interface {
		Value
		Comparable
		Measurable
		Cloneable
		Iterable

		IsEmpty(ctx context.Context) (Boolean, error)
		Clear(ctx context.Context) error
	}

	// List represents a items of values.
	// Generic interface for all items-like structures.
	List interface {
		Collection
		Indexed
		Sortable

		Add(ctx context.Context, value Value) error
		Set(ctx context.Context, idx Int, value Value) error
		Insert(ctx context.Context, idx Int, value Value) error
		Remove(ctx context.Context, value Value) error
		RemoveAt(ctx context.Context, idx Int) (Value, error)

		Find(ctx context.Context, predicate IndexedPredicate) (List, error)
		FindOne(ctx context.Context, predicate IndexedPredicate) (Value, Boolean, error)
		Contains(ctx context.Context, value Value) (Boolean, error)
		IndexOf(ctx context.Context, value Value) (Int, error)
		First(context.Context) (Value, error)
		Last(context.Context) (Value, error)
		Slice(ctx context.Context, start, end Int) (List, error)
		CopyWithCap(ctx context.Context, cap Int) (List, error)

		ForEach(ctx context.Context, predicate IndexedPredicate) error
	}

	// Map represents a dictionary of values.
	// Generic interface for all dictionary-like structures.
	Map interface {
		Collection
		Keyed

		Set(ctx context.Context, key Value, value Value) error
		Remove(ctx context.Context, key Value) error

		ContainsKey(ctx context.Context, key Value) (Boolean, error)
		ContainsValue(ctx context.Context, value Value) (Boolean, error)
		Keys(context.Context) (List, error)
		Values(context.Context) (List, error)
		Find(ctx context.Context, predicate KeyedPredicate) (List, error)
		FindOne(ctx context.Context, predicate KeyedPredicate) (Value, Boolean, error)

		ForEach(ctx context.Context, predicate KeyedPredicate) error
	}
)
