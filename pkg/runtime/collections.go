package runtime

import "context"

type (

	// IndexReadable is an interface for accessing elements by their index in a collection-like structure.
	IndexReadable interface {
		// At retrieves the value at the given index or returns an error if the index is invalid.
		At(ctx context.Context, idx Int) (Value, error)
	}

	// KeyReadable is an interface for accessing elements by their key in a collection-like structure.
	KeyReadable interface {
		// Get retrieves the value associated with the given key or returns an error if the key is not found.
		Get(ctx context.Context, key Value) (Value, error)
	}

	// IndexWritable is an interface for modifying elements by their index in a collection-like structure.
	IndexWritable interface {
		// SetAt method updates the value at the given index or returns an error if the index is invalid.
		SetAt(ctx context.Context, idx Int, value Value) error
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
		// RemoveKey removes the value associated with the specified key from the collection or returns an error if the operation fails.
		RemoveKey(ctx context.Context, key Value) error
	}

	// ValueRemovable is an interface for removing elements by their value in a collection-like structure.
	ValueRemovable interface {
		// Remove method removes the first occurrence of the given value or returns an error if the operation fails.
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

	// Spawnable is an interface for creating new instances of a type.
	// Generic interface for all types that can create new instances of themselves.
	Spawnable[T any] interface {
		// Empty creates a new instance of the type or returns an error if the operation fails.
		// The new instance should be initialized with default values and ready for use.
		Empty(ctx context.Context) (T, error)
	}

	// Collection represents a collection of values.
	// Generic interface for all collection-like structures.
	Collection interface {
		Value
		Containable
		Comparable
		Measurable
		Cloneable
		Iterable

		// Clear removes all elements from the collection or returns an error if the operation fails.
		Clear(ctx context.Context) error
	}

	// List represents a items of values.
	// Generic interface for all items-like structures.
	List interface {
		Collection
		Sortable
		Appendable
		IndexReadable
		IndexWritable
		IndexRemovable
		ValueRemovable
		Spawnable[List]

		// Insert adds a value at the specified index in a collection or returns an error if the index is invalid or operation fails.
		Insert(ctx context.Context, idx Int, value Value) error
		// Swap exchanges the values at two specified indices in a collection or returns an error if indices are invalid or operation fails.
		Swap(ctx context.Context, a, b Int) error
		// Concat concatenates another list to the current list or returns an error if the operation fails.
		Concat(ctx context.Context, other List) error

		// Filter creates a new list containing elements that satisfy the given predicate or returns an error if filtering fails.
		Filter(ctx context.Context, predicate IndexReadablePredicate) (List, error)
		// Find searches for the first element in the list that satisfies the given predicate or returns an error if search fails.
		Find(ctx context.Context, predicate IndexReadablePredicate) (Value, Boolean, error)
		// IndexOf returns the index of the first occurrence of the specified value in the list or -1 if not found.
		IndexOf(ctx context.Context, value Value) (Int, error)
		// First retrieves the first element of the collection or returns an error if the collection is empty or inaccessible.
		First(ctx context.Context) (Value, error)
		// Last retrieves the last element of the collection or returns an error if the collection is empty or inaccessible.
		Last(ctx context.Context) (Value, error)
		// Slice returns a new list containing elements from the specified range or returns an error if range is invalid or inaccessible.
		Slice(ctx context.Context, start, end Int) (List, error)

		// ForEach iterates over each element in the list and applies the given predicate function or returns an error if iteration fails.
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
		Spawnable[Map]

		// Merge merges another map into the current map, combining their key-value pairs or returns an error if merging fails.
		Merge(ctx context.Context, other Map) error

		// ContainsKey checks if the map contains the specified key or returns an error if key check fails.
		ContainsKey(ctx context.Context, key Value) (Boolean, error)
		// Keys returns a list of all keys in the map or an error if retrieval fails.
		Keys(ctx context.Context) (List, error)
		// Values returns a list of all values in the map or an error if retrieval fails.
		Values(ctx context.Context) (List, error)
		// Filter returns a list of values that satisfy the given predicate or an error if filtering fails.
		Filter(ctx context.Context, predicate KeyReadablePredicate) (List, error)
		// Find searches for the first value in the map that satisfies the given predicate or returns an error if search fails.
		Find(ctx context.Context, predicate KeyReadablePredicate) (Value, Boolean, error)

		// ForEach iterates over each key-value pair in the map and applies the given predicate function or returns an error if iteration fails.
		ForEach(ctx context.Context, predicate KeyReadablePredicate) error
	}
)
