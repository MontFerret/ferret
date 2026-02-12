package runtime

import "context"

type (

	// IndexReadable is an interface for accessing elements by their index in a collection-like structure.
	IndexReadable interface {
		// Get retrieves the value at the given index or returns an error if the index is invalid.
		Get(ctx context.Context, idx Int) (Value, error)
	}

	// KeyReadable is an interface for accessing elements by their key in a collection-like structure.
	KeyReadable interface {

		// Get retrieves the value associated with the given key or returns an error if the key is not found.
		Get(ctx context.Context, key Value) (Value, error)
	}

	// IndexWritable is an interface for modifying elements by their index in a collection-like structure.
	IndexWritable interface {
		// Set method updates the value at the given index or returns an error if the index is invalid.
		Set(ctx context.Context, idx Int, value Value) error
	}

	// KeyWritable is an interface for modifying elements by their key in a collection-like structure.
	KeyWritable interface {
		// Set method sets or updates the value associated with the given key or returns an error if the operation fails.
		Set(ctx context.Context, key, value Value) error
	}

	// IndexRemovable is an interface for removing elements by their index in a collection-like structure.
	IndexRemovable interface {
		// RemoveAt method removes the value at the given index and returns it or returns an error if the index is invalid.
		RemoveAt(ctx context.Context, idx Int) (Value, error)
	}

	// KeyRemovable is an interface for removing elements by their key in a collection-like structure.
	KeyRemovable interface {
		// RemoveKey removes the value associated with the specified key from the collection or returns an error if the key is not found.
		RemoveKey(ctx context.Context, key Value) error
	}

	// ValueRemovable is an interface for removing elements by their value in a collection-like structure.
	ValueRemovable interface {
		// Remove method removes the first occurrence of the given value or returns an error if the value is not found.
		Remove(ctx context.Context, value Value) error
	}

	// Appendable is an interface for adding elements to a collection-like structure that supports indexing.
	Appendable interface {
		// Append appends the given value to the end of the collection or returns an error if the operation fails.
		Append(ctx context.Context, val Value) error
	}

	// Containable is an interface for checking the presence of a value in a collection-like structure.
	Containable interface {
		// Contains method checks if the given value exists in the collection and returns a boolean result or an error if the operation fails.
		Contains(ctx context.Context, value Value) (Boolean, error)
	}

	// Predicate is a function type that represents a condition to be evaluated against elements in a collection.
	Predicate = func(ctx context.Context, value, idx Value) (Boolean, error)

	// IndexReadablePredicate is a function type that represents a condition to be evaluated against elements in a collection based on their index.
	IndexReadablePredicate = func(ctx context.Context, value Value, idx Int) (Boolean, error)

	// KeyReadablePredicate is a function type that represents a condition to be evaluated against elements in a collection based on their key.
	KeyReadablePredicate = func(ctx context.Context, value, key Value) (Boolean, error)

	// Collection represents a collection of values.
	// Generic interface for all collection-like structures.
	Collection interface {
		Value
		Containable
		Comparable
		Measurable
		Cloneable
		Iterable

		Clear(ctx context.Context) error
	}

	// List represents a items of values.
	// Generic interface for all items-like structures.
	List interface {
		Collection
		Appendable
		IndexReadable
		IndexWritable
		IndexRemovable
		ValueRemovable

		Insert(ctx context.Context, idx Int, value Value) error
		Swap(ctx context.Context, a, b Int) error

		Find(ctx context.Context, predicate IndexReadablePredicate) (List, error)
		FindOne(ctx context.Context, predicate IndexReadablePredicate) (Value, Boolean, error)
		IndexOf(ctx context.Context, value Value) (Int, error)
		First(ctx context.Context) (Value, error)
		Last(ctx context.Context) (Value, error)
		Slice(ctx context.Context, start, end Int) (List, error)

		ForEach(ctx context.Context, predicate IndexReadablePredicate) error
	}

	// Map represents a dictionary of values.
	// Generic interface for all dictionary-like structures.
	Map interface {
		Collection
		KeyReadable
		KeyWritable
		KeyRemovable
		ValueRemovable

		ContainsKey(ctx context.Context, key Value) (Boolean, error)
		Keys(ctx context.Context) (List, error)
		Values(ctx context.Context) (List, error)
		Find(ctx context.Context, predicate KeyReadablePredicate) (List, error)
		FindOne(ctx context.Context, predicate KeyReadablePredicate) (Value, Boolean, error)

		ForEach(ctx context.Context, predicate KeyReadablePredicate) error
	}
)
