package runtime

type (

	// IndexReadable is an interface for accessing elements by their index in a collection-like structure.
	// The Get method retrieves the value at the given index or returns an error if the index is invalid.
	IndexReadable interface {
		Get(ctx Context, idx Int) (Value, error)
	}

	// KeyReadable is an interface for accessing elements by their key in a collection-like structure.
	// The Get method retrieves the value associated with the given key or returns an error if the key is not found.
	KeyReadable interface {
		Get(ctx Context, key Value) (Value, error)
	}

	// IndexWritable is an interface for modifying elements by their index in a collection-like structure.
	// The Set method updates the value at the given index or returns an error if the index is invalid.
	IndexWritable interface {
		Set(ctx Context, idx Int, value Value) error
	}

	// KeyWritable is an interface for modifying elements by their key in a collection-like structure.
	// The Set method updates the value associated with the given key or returns an error if the key is not found.
	KeyWritable interface {
		Set(ctx Context, key, value Value) error
	}

	// IndexAppendable is an interface for adding elements to a collection-like structure that supports indexing.
	// The Append method appends the given value to the end of the collection or returns an error if the operation fails.
	IndexAppendable interface {
		Append(ctx Context, val Value) error
	}

	// IndexRemovable is an interface for removing elements by their index in a collection-like structure.
	// The RemoveAt method removes the value at the given index and returns it or returns an error if the index is invalid.
	IndexRemovable interface {
		RemoveAt(ctx Context, idx Int) (Value, error)
	}

	// KeyRemovable is an interface for removing elements by their key in a collection-like structure.
	// The RemoveKey method removes the value associated with the given key or returns an error if the key is not found.
	KeyRemovable interface {
		RemoveKey(ctx Context, key Value) error
	}

	// ValueRemovable is an interface for removing elements by their value in a collection-like structure.
	// The RemoveValue method removes the first occurrence of the given value or returns an error if the value is not found.
	ValueRemovable interface {
		RemoveValue(ctx Context, value Value) error
	}

	// Predicate is a function type that represents a condition to be evaluated against elements in a collection.
	Predicate = func(ctx Context, value, idx Value) (Boolean, error)

	// IndexReadablePredicate is a function type that represents a condition to be evaluated against elements in a collection based on their index.
	IndexReadablePredicate = func(ctx Context, value Value, idx Int) (Boolean, error)

	// KeyReadablePredicate is a function type that represents a condition to be evaluated against elements in a collection based on their key.
	KeyReadablePredicate = func(ctx Context, value, key Value) (Boolean, error)

	// Collection represents a collection of values.
	// Generic interface for all collection-like structures.
	Collection interface {
		Value
		Comparable
		Measurable
		Cloneable
		Iterable

		Clear(ctx Context) error
	}

	// List represents a items of values.
	// Generic interface for all items-like structures.
	List interface {
		Collection
		IndexReadable
		IndexWritable
		IndexAppendable
		IndexRemovable
		ValueRemovable

		Insert(ctx Context, idx Int, value Value) error
		Swap(ctx Context, a, b Int) error

		Find(ctx Context, predicate IndexReadablePredicate) (List, error)
		FindOne(ctx Context, predicate IndexReadablePredicate) (Value, Boolean, error)
		IndexOf(ctx Context, value Value) (Int, error)
		First(ctx Context) (Value, error)
		Last(ctx Context) (Value, error)
		Slice(ctx Context, start, end Int) (List, error)

		ForEach(ctx Context, predicate IndexReadablePredicate) error
	}

	// Map represents a dictionary of values.
	// Generic interface for all dictionary-like structures.
	Map interface {
		Collection
		KeyReadable
		KeyWritable
		KeyRemovable
		ValueRemovable

		ContainsKey(ctx Context, key Value) (Boolean, error)
		ContainsValue(ctx Context, value Value) (Boolean, error)
		Keys(ctx Context) (List, error)
		Values(ctx Context) (List, error)
		Find(ctx Context, predicate KeyReadablePredicate) (List, error)
		FindOne(ctx Context, predicate KeyReadablePredicate) (Value, Boolean, error)

		ForEach(ctx Context, predicate KeyReadablePredicate) error
	}
)
